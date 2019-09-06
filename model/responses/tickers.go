package responses

type Tickers struct {
	Code      int64    `json:"code"`
	Timestamp int64    `json:"timestamp"`
	Tickers   []Ticker `json:"tickers"`
}

type Ticker struct {
	Symbol      string `json:"symbol"`
	AskPrice    string `json:"askPrice"`
	AskQty      string `json:"askQty"`
	BidPrice    string `json:"bidPrice"`
	BidQty      string `json:"bidQty"`
	LastPrice   string `json:"lastPrice"`
	LastQty     string `json:"lastQty"`
	Volume      string `json:"volume"`
	QuoteVolume string `json:"quoteVolume"`
	OpenTime    int64  `json:"openTime"`
	CloseTime   int64  `json:"closeTime"`
}
