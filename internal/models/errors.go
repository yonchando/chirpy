package models

type ErrorBody struct {
	Message string      `json:"message"`
	Errors  interface{} `json:"errors"`
}
