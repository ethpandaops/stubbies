package execution

type RequestParamsForkchoiceUpdatedV1 struct {
	HeadBlockHash string `json:"headBlockHash"`
}

type RequestParamsNewPayloadV1 struct {
	ParentHash    string   `json:"parentHash"`
	FeeRecipient  string   `json:"feeRecipient"`
	StateRoot     string   `json:"stateRoot"`
	ReceiptsRoot  string   `json:"receiptsRoot"`
	LogsBloom     string   `json:"logsBloom"`
	Random        string   `json:"prevRandao"`
	BlockNumber   string   `json:"blockNumber"`
	GasLimit      string   `json:"gasLimit"`
	GasUsed       string   `json:"gasUsed"`
	Timestamp     string   `json:"timestamp"`
	ExtraData     string   `json:"extraData"`
	BaseFeePerGas string   `json:"baseFeePerGas"`
	BlockHash     string   `json:"blockHash"`
	Transactions  []string `json:"transactions"`
}

type RequestParamsExchangeCapabilities []string
