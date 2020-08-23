package models

// ResponseError formats the error sent to the client
type ResponseError struct {
	Status 	int
	Message		string
}