package cmd

type Response struct {
	Error *ErrorResponse `json:"error,omitempty"`
	Data  []byte         `json:"data"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}
