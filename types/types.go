package types

import "github.com/gocql/gocql"

type ExampleEntity struct {
	ID      gocql.UUID `json:"id"`
	Message string     `json:"message"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}
