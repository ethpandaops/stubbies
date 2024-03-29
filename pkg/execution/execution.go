package execution

import (
	"context"
	"encoding/json"
	"errors"
	"math/big"

	"github.com/sirupsen/logrus"
)

var ErrUnsupportedGetBlockQuery = errors.New("unsupported get block query")

type Handler struct {
	log logrus.FieldLogger
	Cfg Config

	storage *Storage
}

// NewHandler returns a new Handler instance.
func NewHandler(log logrus.FieldLogger, conf *Config) *Handler {
	if err := conf.Validate(); err != nil {
		log.Fatalf("invalid config: %s", err)
	}

	return &Handler{
		log:     log.WithField("module", "api/execution"),
		Cfg:     *conf,
		storage: newStorage(log.WithField("module", "api/execution/storage")),
	}
}

func (h *Handler) Start(ctx context.Context) {
	h.storage.Start(ctx)
}

func (h *Handler) Request(ctx context.Context, id int, method string, params []*json.RawMessage) (*Response, error) {
	resp := &Response{
		ID:      id,
		JSONRPC: "2.0",
		Result:  false,
	}

	switch method {
	case "eth_syncing":
		resp.Result = false
	case "engine_exchangeTransitionConfigurationV1":
		resp.Result = ResultExchangeTransitionConfigurationV1{
			TerminalTotalDifficulty: h.Cfg.TerminalTotalDifficulty,
			TerminalBlockHash:       h.Cfg.TerminalBlockHash,
			TerminalBlockNumber:     h.Cfg.TerminalBlockNumber,
		}
	case "engine_forkchoiceUpdatedV1":
		result, err := h.forkChoiceUpdated(params)
		if err != nil {
			return nil, err
		}

		resp.Result = result
	case "engine_forkchoiceUpdatedV2":
		result, err := h.forkChoiceUpdated(params)
		if err != nil {
			return nil, err
		}

		resp.Result = result
	case "engine_newPayloadV1":
		result, err := h.newPayload(params)
		if err != nil {
			return nil, err
		}

		resp.Result = result
	case "engine_newPayloadV2":
		result, err := h.newPayload(params)
		if err != nil {
			return nil, err
		}

		resp.Result = result
	case "engine_newPayloadV3":
		result, err := h.newPayload(params)
		if err != nil {
			return nil, err
		}

		resp.Result = result
	case "eth_getBlockByHash":
		result, err := h.getBlockByHash(params)
		if err != nil && err != ErrUnsupportedGetBlockQuery {
			return nil, err
		}

		resp.Result = result
	case "eth_getBlockByNumber":
		result, err := h.getBlockByNumber(params)
		if err != nil && err != ErrUnsupportedGetBlockQuery {
			return nil, err
		}

		resp.Result = result
	case "eth_chainId":
		resp.Result = ResultChainID(h.Cfg.ChainID)
	case "engine_exchangeCapabilities":
		if len(params) < 1 || params[0] == nil {
			return nil, errors.New("missing params")
		}

		var payload RequestParamsExchangeCapabilities

		err := json.Unmarshal([]byte(*params[0]), &payload)
		if err != nil {
			return nil, err
		}

		resp.Result = ResultexchangeCapabilities(payload)
	case "eth_call":
	default:
		h.log.WithField("method", method).Warn("unsupported method")
	}

	return resp, nil
}

func (h *Handler) forkChoiceUpdated(params []*json.RawMessage) (interface{}, error) {
	if len(params) < 1 || params[0] == nil {
		return nil, errors.New("missing params")
	}

	var forkchoiceState RequestParamsForkchoiceUpdatedV1

	err := json.Unmarshal([]byte(*params[0]), &forkchoiceState)
	if err != nil {
		return nil, err
	}

	return ResultForkchoiceUpdatedV1{
		PayloadStatus: ResultForkchoiceUpdatedV1PayloadStatus{
			Status:          "VALID",
			LatestValidHash: forkchoiceState.HeadBlockHash,
			ValidationError: "",
		},
		PayloadID: "0xa247243752eb10b4",
	}, nil
}

func (h *Handler) newPayload(params []*json.RawMessage) (interface{}, error) {
	if len(params) < 1 || params[0] == nil {
		return nil, errors.New("missing params")
	}

	var payload RequestParamsNewPayloadV1

	err := json.Unmarshal([]byte(*params[0]), &payload)
	if err != nil {
		return nil, err
	}

	h.storage.AddBlock(&payload, params[0])

	return ResultNewPayloadV1{
		Status:          "VALID",
		LatestValidHash: payload.BlockHash,
		ValidationError: "",
	}, nil
}

func (h *Handler) getBlockByHash(params []*json.RawMessage) (interface{}, error) {
	if len(params) < 1 || params[0] == nil {
		return nil, errors.New("missing params")
	}

	var query string

	err := json.Unmarshal([]byte(*params[0]), &query)
	if err != nil {
		return nil, err
	}

	if query == "latest" {
		block := h.storage.GetLatestBlock()
		if block != nil {
			return block.GetResult(), nil
		}
	} else {
		block := h.storage.GetBlockByHash(query)
		if block != nil {
			return block.GetResult(), nil
		}
	}

	return "{}", ErrUnsupportedGetBlockQuery
}

func (h *Handler) getBlockByNumber(params []*json.RawMessage) (interface{}, error) {
	if len(params) < 1 || params[0] == nil {
		return nil, errors.New("missing params")
	}

	var query string

	err := json.Unmarshal([]byte(*params[0]), &query)
	if err != nil {
		return nil, err
	}

	if query == "latest" {
		block := h.storage.GetLatestBlock()
		if block != nil {
			return block.GetResult(), nil
		}
	} else {
		number := new(big.Int)
		number.SetString(query[2:], 16)
		block := h.storage.GetBlockByNumber(number)
		if block != nil {
			return block.GetResult(), nil
		}
	}

	return nil, ErrUnsupportedGetBlockQuery
}
