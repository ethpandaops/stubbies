package execution

import (
	"encoding/json"
	"math/big"
)

type Block struct {
	Number *big.Int
	raw    *json.RawMessage
}

func (b *Block) UpdateToLatest(payload RequestParamsNewPayloadV1, raw *json.RawMessage) {
	num := new(big.Int)
	num.SetString(payload.BlockNumber[2:], 16)

	// replace if the new block is higher
	if num.Cmp(b.Number) > 0 {
		b.Number = num
		b.raw = raw
	}
}
