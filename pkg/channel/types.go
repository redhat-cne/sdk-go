package channel

//Status specifies status of the event
type Status int

const (
	// SUCCEED if the event is posted successfully
	SUCCEED Status = 1
	//FAILED if the event  failed to post
	FAILED Status = 2
	//NEW if the event is new for the consumer
	NEW Status = 0
	//DELETE if the event is to delete
	DELETE Status = -1
)

//Type ... specifies type of the event
type Type int

const (
	// LISTENER  the type to create listener
	LISTENER Type = 1
	//STATUS  the  types is check status
	STATUS Type = 2
	//SENDER  the  type is to create sender
	SENDER Type = 0
	//EVENT  the type is an event
	EVENT Type = 3
)
