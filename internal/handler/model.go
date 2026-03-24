package handler

type PackRequest struct {
	Input string `json:"input"`
}

type PackResponse struct {
	Result string `json:"result"`
}

type UnpackRequest struct {
	Input string `json:"input"`
}

type UnpackResponse struct {
	RequestID string `json:"request_id"`
	Result    string `json:"result"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}
