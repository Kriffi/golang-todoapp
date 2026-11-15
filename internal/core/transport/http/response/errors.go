package core_http_response

type ErrorResponseBody struct {
	Message string `json:"message"`
	Error   string `json:"error"`
}
