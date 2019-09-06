package requests

type Order struct {
	Symbol    string `json:"symbol"`
	Side      string `json:"side"`
	Type      string `json:"type"`
	Price     float64  `json:"price"`
	StopPrice float64  `json:"stopPrice"`
	Quantity  float64  `json:"quantity"`
	Timestamp int64 `json:"ts"`
}