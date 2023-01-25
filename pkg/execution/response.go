package execution

import "encoding/json"

type Response struct {
	ID      int         `json:"id"`
	JSONRPC string      `json:"jsonrpc"`
	Result  interface{} `json:"result"`
}

type ResultDefault bool

type ResultExchangeTransitionConfigurationV1 struct {
	TerminalTotalDifficulty string `json:"terminalTotalDifficulty"`
	TerminalBlockHash       string `json:"terminalBlockHash"`
	TerminalBlockNumber     string `json:"terminalBlockNumber"`
}

type ResultForkchoiceUpdatedV1 struct {
	PayloadStatus ResultForkchoiceUpdatedV1PayloadStatus `json:"payloadStatus"`
	PayloadID     string                                 `json:"payloadId"`
}

type ResultForkchoiceUpdatedV1PayloadStatus struct {
	Status          string `json:"status"`
	LatestValidHash string `json:"latestValidHash"`
	ValidationError string `json:"validationError"`
}

type ResultNewPayloadV1 struct {
	Status          string `json:"status"`
	LatestValidHash string `json:"latestValidHash"`
	ValidationError string `json:"validationError"`
}

type ResultGetBlockByNumber *json.RawMessage

type ResultChainID string
