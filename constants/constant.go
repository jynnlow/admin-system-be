package constants

type Response struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

const (
	SUCCESS = "SUCCESS"
	FAIL    = "FAIL"
)
