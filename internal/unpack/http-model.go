package unpack

type PackHTTPRequest struct {
	Input string `json:"input"` 
}

type PackHTTPResponse struct {
	Result string `json:"result"`
}

type UnpackHTTPRequest struct {
	Input string `json:"input"`
}

type UnpackHTTPResponse struct {
	RequestID string `json:"request_id"`
	Result    string `json:"result"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}