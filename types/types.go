package types

import (
	gjson "github.com/goccy/go-json"
)

type TraceData struct {
	TxHash string           `json:"txHash"`
	Result gjson.RawMessage `json:"result"`
}

type TransactionData struct {
	TransactionHash string                 `json:"transactionHash"`
	Data            map[string]interface{} `json:"-"`
}
