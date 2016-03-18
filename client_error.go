package main

const ClientErrorInvalidMessage = "invalid_message"
const ClientErrorConnection = "connection_error"
const ClientErrorInvalidReceivers = "invalid_receivers"

type ClientError struct {
	problem string
}

func newClientInvalidCommandError() *ClientError {
	return &ClientError{problem: ClientErrorInvalidMessage}
}

func newClientConnectionError() *ClientError {
	return &ClientError{problem: ClientErrorConnection}
}

func newClientInvalidReceiversError() *ClientError {
	return &ClientError{problem: ClientErrorInvalidReceivers}
}

func (e *ClientError) InvalidCommand() bool {
	return e.problem == ClientErrorInvalidMessage
}

func (e *ClientError) ConnectionError() bool {
	return e.problem == ClientErrorConnection
}

func (e *ClientError) InvalidReceivers() bool {
	return e.problem == ClientErrorInvalidReceivers
}

// Implementing error interface
func (e *ClientError) Error() string {
	return e.problem
}
