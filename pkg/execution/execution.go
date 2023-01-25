package execution

import (
	"context"
	"encoding/json"
	"errors"
	"math/big"

	"github.com/sirupsen/logrus"
)

type Handler struct {
	log logrus.FieldLogger
	Cfg Config

	latestBlock Block
}

// NewHandler returns a new Handler instance.
func NewHandler(log logrus.FieldLogger, conf *Config) *Handler {
	if err := conf.Validate(); err != nil {
		log.Fatalf("invalid config: %s", err)
	}

	return &Handler{
		log: log.WithField("module", "api/execution"),
		Cfg: *conf,
		latestBlock: Block{
			Number: big.NewInt(0),
			raw:    nil,
		},
	}
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
		if len(params) < 1 || params[0] == nil {
			return nil, errors.New("missing params")
		}

		var forkchoiceState RequestParamsForkchoiceUpdatedV1

		err := json.Unmarshal([]byte(*params[0]), &forkchoiceState)
		if err != nil {
			return nil, err
		}

		resp.Result = ResultForkchoiceUpdatedV1{
			PayloadStatus: ResultForkchoiceUpdatedV1PayloadStatus{
				Status:          "VALID",
				LatestValidHash: forkchoiceState.HeadBlockHash,
				ValidationError: "",
			},
			PayloadID: "0xa247243752eb10b4",
		}
	case "engine_newPayloadV1":
		if len(params) < 1 || params[0] == nil {
			return nil, errors.New("missing params")
		}

		var payload RequestParamsNewPayloadV1

		err := json.Unmarshal([]byte(*params[0]), &payload)
		if err != nil {
			return nil, err
		}

		h.latestBlock.UpdateToLatest(payload, params[0])

		resp.Result = ResultNewPayloadV1{
			Status:          "VALID",
			LatestValidHash: payload.BlockHash,
			ValidationError: "",
		}
	case "eth_getBlockByNumber":
		if len(params) < 1 || params[0] == nil {
			return nil, errors.New("missing params")
		}

		var query string

		err := json.Unmarshal([]byte(*params[0]), &query)
		if err != nil {
			return nil, err
		}

		if query != "latest" {
			resp.Result = ResultGetBlockByNumber(h.latestBlock.raw)
		}
	case "eth_chainId":
		resp.Result = ResultChainID(h.Cfg.ChainID)
	default:
	}

	return resp, nil
}
