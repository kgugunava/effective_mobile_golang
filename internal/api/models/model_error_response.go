package models

type ErrorResponse struct {
	Error ErrorResponseError `json:"error"`
}