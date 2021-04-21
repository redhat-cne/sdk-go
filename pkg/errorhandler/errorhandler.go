package errorhandler

import "fmt"

// ErrorHandler  ... custom error handler  interface
type ErrorHandler interface {
	Error() string
}

// ReceiverNotFoundError  ...
type ReceiverNotFoundError struct {
	Name string
	Desc string
}

// Error  ...
func (r ReceiverNotFoundError) Error() string {
	return fmt.Sprintf("receiver %s not found", r.Name)
}

// ReceiverError  ...
type ReceiverError struct {
	Name string
	Desc string
}

// Error ...
func (r ReceiverError) Error() string {
	return fmt.Sprintf("receiver %s error %s", r.Name, r.Desc)
}

// SenderError ...
type SenderError struct {
	Name string
	Desc string
}

// Error ...
func (sr SenderError) Error() string {
	return fmt.Sprintf("sender %s error %s", sr.Name, sr.Desc)
}

// SenderNotFoundError  ... sender not found custom error
type SenderNotFoundError struct {
	Name string
	Desc string
}

//Error  ...
func (s SenderNotFoundError) Error() string {
	return fmt.Sprintf("sender %s not found", s.Name)
}

// AMQPConnectionError  ... custom amqp connection error
type AMQPConnectionError struct {
	Desc string
}

// Error ...
func (a AMQPConnectionError) Error() string {
	return fmt.Sprintf("amqp connection error %s", a.Desc)
}

// CloudEventsClientError ... custom client error
type CloudEventsClientError struct {
	Desc string
}

// Error ...
func (c CloudEventsClientError) Error() string {
	return fmt.Sprintf(" cloud events client error %s", c.Desc)
}
