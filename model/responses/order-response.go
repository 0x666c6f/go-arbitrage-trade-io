package responses

type OrderResponse struct {
	Code      int `json:"code"`
	Timestamp int `json:"timestamp"`
	Order     Order `json:"order"`
}

type Order struct {
	OrderID         string `json:"orderId"`
	Total           string `json:"total"`
	OrderType       string `json:"orderType"`
	Commission      string `json:"commission"`
	CreatedAt       string `json:"createdAt"`
	UnitsFilled     string `json:"unitsFilled"`
	IsPending       bool   `json:"isPending"`
	Status          string `json:"status"`
	Type            string `json:"type"`
	RequestedAmount string `json:"requestedAmount"`
	BaseAmount      string `json:"baseAmount"`
	QuoteAmount     string `json:"quoteAmount"`
	Price           string `json:"price"`
	IsLimit         bool   `json:"isLimit"`
	LoanRate        string `json:"loanRate"`
	RateStop        string `json:"rateStop"`
	Instrument      string `json:"instrument"`
	RequestedPrice  string `json:"requestedPrice"`
	RemainingAmount string `json:"remainingAmount"`
}

