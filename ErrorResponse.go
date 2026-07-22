package faire_go

// ErrorResponse stores general API error response
type ErrorResponse struct {
	StatusCode   int           `json:"status_code,omitempty"`
	StatusType   string        `json:"status_type,omitempty"`
	Message      string        `json:"message,omitempty"`
	Field        string        `json:"field,omitempty"`
	EntityToken  string        `json:"entity_token,omitempty"`
	EntityTokens []interface{} `json:"entity_tokens,omitempty"`
}
