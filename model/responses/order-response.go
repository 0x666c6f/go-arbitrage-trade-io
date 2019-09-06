package responses

type OrderResponse struct {
	Code      int `json:"code"`
	Timestamp int `json:"timestamp"`
	Order     Order `json:"order"`
}

type Order struct {
	OrderID         string `json:"orderId"`
	Total           float64 `json:"total"`
	OrderType       string `json:"orderType"`
	Commission      float64 `json:"commission"`
	CreatedAt       string `json:"createdAt"`
	UnitsFilled     float64 `json:"unitsFilled"`
	IsPending       bool   `json:"isPending"`
	Status          string `json:"status"`
	Type            string `json:"type"`
	RequestedAmount float64 `json:"requestedAmount"`
	BaseAmount      float64 `json:"baseAmount"`
	QuoteAmount     float64 `json:"quoteAmount"`
	Price           float64 `json:"price"`
	IsLimit         bool   `json:"isLimit"`
	LoanRate        float64 `json:"loanRate"`
	RateStop        float64 `json:"rateStop"`
	Instrument      string `json:"instrument"`
	RequestedPrice  float64 `json:"requestedPrice"`
	RemainingAmount float64 `json:"remainingAmount"`
}

