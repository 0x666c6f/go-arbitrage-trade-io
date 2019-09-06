package responses

type Infos struct {
	Code      int64    `json:"code"`
	Timestamp int64    `json:"timestamp"`
	Symbols   []Symbol `json:"symbols"`
}

type Symbol struct {
	Symbol              string `json:"symbol"`
	Status              string `json:"status"`
	BaseAsset           string `json:"baseAsset"`
	BaseAssetPrecision  int64  `json:"baseAssetPrecision"`
	QuoteAsset          string `json:"quoteAsset"`
	QuoteAssetPrecision int64  `json:"quoteAssetPrecision"`
}
