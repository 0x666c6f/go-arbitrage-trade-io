package responses

type Balances struct {
	Code      int     `json:"code"`
	Timestamp int     `json:"timestamp"`
	Balances  []Balance `json:"balances"`
}

type Balance struct {
	Asset     string `json:"asset"`
	Available float64 `json:"available"`
	Locked    float64 `json:"locked"`
}
