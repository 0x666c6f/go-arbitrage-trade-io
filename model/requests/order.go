package requests

type Order struct {
	Symbol    string `json:"symbol"`
	Side      string `json:"side"`
	Type      string `json:"type"`
	Price     int64  `json:"price"`
	StopPrice int64  `json:"stopPrice"`
	Quantity  float64  `json:"quantity"`
	Timestamp string `json:"ts"`
}