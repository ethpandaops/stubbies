package execution

import (
	"encoding/json"
	"math/big"
)

type Block struct {
	Number  *big.Int
	raw     *json.RawMessage
	payload *RequestParamsNewPayloadV1
}

func (b *Block) UpdateToLatest(payload *RequestParamsNewPayloadV1, raw *json.RawMessage) {
	num := new(big.Int)
	num.SetString(payload.BlockNumber[2:], 16)

	// replace if the new block is higher
	if num.Cmp(b.Number) > 0 {
		b.Number = num
		b.raw = raw
		b.payload = payload
	}
}

func (b *Block) GetResult() *ResultGetBlock {
	if b.payload == nil {
		return nil
	}

	return &ResultGetBlock{
		Number:       b.payload.BlockNumber,
		Hash:         b.payload.BlockHash,
		ParentHash:   b.payload.ParentHash,
		LogsBloom:    b.payload.LogsBloom,
		StateRoot:    b.payload.StateRoot,
		ReceiptsRoot: b.payload.ReceiptsRoot,
		ExtraData:    b.payload.ExtraData,
		GasLimit:     b.payload.GasLimit,
		GasUsed:      b.payload.GasUsed,
		Timestamp:    b.payload.Timestamp,
		Transactions: b.payload.Transactions,
	}
}
