package responses

type Infos struct {
	Code      int    `json:"code"`
	Timestamp int    `json:"timestamp"`
	Symbols   []Symbol `json:"symbols"`
}

type Symbol struct {
	Symbol              string `json:"symbol"`
	Status              string `json:"status"`
	BaseAsset           string `json:"baseAsset"`
	BaseAssetPrecision  int  `json:"baseAssetPrecision"`
	QuoteAsset          string `json:"quoteAsset"`
	QuoteAssetPrecision int  `json:"quoteAssetPrecision"`
}
