package api

import "encoding/json"

type JSONRequestBody struct {
	ID     int                `json:"id"`
	Method string             `json:"method"`
	Params []*json.RawMessage `json:"params"`
}
