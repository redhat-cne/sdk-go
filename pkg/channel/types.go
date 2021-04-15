package channel

//Status specifies status of the event
type Status int

const (
	//NEW if the event is new for the consumer
	NEW Status = iota
	// SUCCEED if the event is posted successfully
	SUCCEED
	//DELETE if the event is to delete
	DELETE
	//FAILED if the event  failed to post
	FAILED
)

//Type ... specifies type of the event
type Type int

const (
	// LISTENER  the type to create listener
	LISTENER Type = iota
	//SENDER  the  type is to create sender
	SENDER
	//EVENT  the type is an event
	EVENT
)
