package responses

type Balances struct {
	Code      int64     `json:"code"`
	Timestamp int64     `json:"timestamp"`
	Balances  []Balance `json:"balances"`
}

type Balance struct {
	Asset     string `json:"asset"`
	Available string `json:"available"`
	Locked    string `json:"locked"`
}
