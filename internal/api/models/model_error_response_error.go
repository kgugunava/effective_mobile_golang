package models

type ErrorResponseError struct {
	Code string `json:"code"`
	Message string `json:"message"`
}