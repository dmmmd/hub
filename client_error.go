package main

const ClientErrorInvalidMessage = "invalid_message"
const ClientErrorConnectionLost = "connection_lost"
const ClientErrorInvalidReceivers = "invalid_receivers"

type ClientError struct {
	problem string
}

func NewClientInvalidMessageError() *ClientError {
	return &ClientError{problem: ClientErrorInvalidMessage}
}

func NewClientConnectionLostError() *ClientError {
	return &ClientError{problem: ClientErrorConnectionLost}
}

func NewClientInvalidReceivers() *ClientError {
	return &ClientError{problem: ClientErrorInvalidReceivers}
}

func (e *ClientError) InvalidMessage() bool {
	return e.problem == ClientErrorInvalidMessage
}

func (e *ClientError) ConnectionLost() bool {
	return e.problem == ClientErrorConnectionLost
}

func (e *ClientError) InvalidReceivers() bool {
	return e.problem == ClientErrorInvalidReceivers
}

// Implementing error interface
func (e *ClientError) Error() string {
	return e.problem
}
