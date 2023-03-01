package execution

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

type ResultChainID string

type ResultexchangeCapabilities []string

type ResultGetBlock struct {
	Number       string   `json:"number"`
	Hash         string   `json:"hash"`
	ParentHash   string   `json:"parentHash"`
	LogsBloom    string   `json:"logsBloom"`
	StateRoot    string   `json:"stateRoot"`
	ReceiptsRoot string   `json:"receiptsRoot"`
	ExtraData    string   `json:"extraData"`
	GasLimit     string   `json:"gasLimit"`
	GasUsed      string   `json:"gasUsed"`
	Timestamp    string   `json:"timestamp"`
	Transactions []string `json:"transactions"`
}
