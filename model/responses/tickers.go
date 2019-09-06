package responses

type Tickers struct {
	Code      int    `json:"code"`
	Timestamp int    `json:"timestamp"`
	Tickers   []Ticker `json:"tickers"`
}

type Ticker struct {
	Symbol      string `json:"symbol"`
	AskPrice    float64 `json:"askPrice"`
	AskQty      float64 `json:"askQty"`
	BidPrice    float64 `json:"bidPrice"`
	BidQty      float64 `json:"bidQty"`
	LastPrice   float64 `json:"lastPrice"`
	LastQty     float64 `json:"lastQty"`
	Volume      float64 `json:"volume"`
	QuoteVolume float64 `json:"quoteVolume"`
	OpenTime    int  `json:"openTime"`
	CloseTime   int  `json:"closeTime"`
}
