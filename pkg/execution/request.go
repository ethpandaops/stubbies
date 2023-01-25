package execution

type RequestParamsForkchoiceUpdatedV1 struct {
	HeadBlockHash string `json:"headBlockHash"`
}

type RequestParamsNewPayloadV1 struct {
	BlockHash   string `json:"blockHash"`
	BlockNumber string `json:"blockNumber"`
}
