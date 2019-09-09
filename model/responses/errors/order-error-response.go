package errors

type OrderErrorResponse struct {
	Errors []OrderError `json:"errors"`
}

type OrderError struct {
	Code string `json:"code"`
	Message string `json:"message"`
	Value float64 `json:"value"`
	Max float64 `json:"max"`
}