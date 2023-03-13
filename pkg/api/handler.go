package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	exec "github.com/ethpandaops/stubbies/pkg/execution"
	"github.com/julienschmidt/httprouter"
	"github.com/sirupsen/logrus"
)

// Handler is an API handler that is responsible for negotiating with a HTTP api.
type Handler struct {
	log logrus.FieldLogger

	execution *exec.Handler

	metrics Metrics
}

func NewHandler(log logrus.FieldLogger, conf *exec.Config) *Handler {
	return &Handler{
		log: log.WithField("module", "api"),

		execution: exec.NewHandler(log, conf),

		metrics: NewMetrics("http"),
	}
}

func (h *Handler) Register(ctx context.Context, router *httprouter.Router) error {
	router.POST("/", h.wrappedHandler(h.handleExecution))

	return nil
}

func (h *Handler) Start(ctx context.Context) {
	h.execution.Start(ctx)
}

func deriveRegisteredPath(request *http.Request, ps httprouter.Params) string {
	registeredPath := request.URL.Path
	for _, param := range ps {
		registeredPath = strings.Replace(registeredPath, param.Value, fmt.Sprintf(":%s", param.Key), 1)
	}

	return registeredPath
}

func (h *Handler) wrappedHandler(handler func(ctx context.Context, r *http.Request, p httprouter.Params, contentType ContentType, body *JSONRequestBody) (*HTTPResponse, error)) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		start := time.Now()

		contentType := NewContentTypeFromRequest(r)
		ctx := r.Context()
		registeredPath := deriveRegisteredPath(r, p)

		h.log.WithFields(logrus.Fields{
			"method":       r.Method,
			"path":         r.URL.Path,
			"content_type": contentType,
			"accept":       r.Header.Get("Accept"),
		}).Debug("handling http request")

		executionMethod := "unknown"

		h.metrics.ObserveRequest(r.Method, registeredPath, executionMethod)

		responseStatusCode := http.StatusInternalServerError

		var err error

		defer func() {
			h.metrics.ObserveResponse(r.Method, registeredPath, fmt.Sprintf("%v", responseStatusCode), contentType.String(), executionMethod, time.Since(start))
		}()

		decoder := json.NewDecoder(r.Body)

		var body JSONRequestBody

		err = decoder.Decode(&body)
		if err != nil {
			responseStatusCode = http.StatusBadRequest
			if writeErr := WriteErrorResponse(w, err.Error(), responseStatusCode); writeErr != nil {
				h.log.WithError(writeErr).Error("Failed to decode request body")
			}

			return
		}

		if body.Method == "" {
			responseStatusCode = http.StatusBadRequest
			if writeErr := WriteErrorResponse(w, "Missing method", responseStatusCode); writeErr != nil {
				h.log.WithError(writeErr).Error("Request body missing method")
			}

			return
		}

		executionMethod = body.Method

		response, err := handler(ctx, r, p, contentType, &body)
		if err != nil {
			if response != nil && response.StatusCode != 0 {
				responseStatusCode = response.StatusCode
			}

			if writeErr := WriteErrorResponse(w, err.Error(), responseStatusCode); writeErr != nil {
				h.log.WithError(writeErr).Error("Failed to write error response")
			}

			return
		}

		responseStatusCode = response.StatusCode

		data, err := response.MarshalAs(contentType)
		if err != nil {
			responseStatusCode = http.StatusInternalServerError
			if writeErr := WriteErrorResponse(w, err.Error(), responseStatusCode); writeErr != nil {
				h.log.WithError(writeErr).Error("Failed to write error response")
			}

			return
		}

		for header, value := range response.Headers {
			w.Header().Set(header, value)
		}

		if err := WriteContentAwareResponse(w, data, contentType); err != nil {
			h.log.WithError(err).Error("Failed to write response")
		}
	}
}

func (h *Handler) handleExecution(ctx context.Context, r *http.Request, p httprouter.Params, contentType ContentType, body *JSONRequestBody) (*HTTPResponse, error) {
	if err := ValidateContentType(contentType, []ContentType{ContentTypeJSON}); err != nil {
		return NewUnsupportedMediaTypeResponse(nil), err
	}

	var parms []string

	for _, param := range body.Params {
		if param != nil {
			parms = append(parms, string(*param))
		}
	}

	h.log.WithFields(logrus.Fields{
		"method": body.Method,
		"params": parms,
		"id":     body.ID,
	}).Debug("handling execution request")

	resp, err := h.execution.Request(ctx, body.ID, body.Method, body.Params)
	if err != nil {
		return NewInternalServerErrorResponse(nil), err
	}

	var rsp = NewSuccessResponse(ContentTypeResolvers{
		ContentTypeJSON: func() ([]byte, error) {
			return json.Marshal(resp)
		},
	})

	rsp.SetCacheControl("public, s-max-age=30")

	return rsp, nil
}
