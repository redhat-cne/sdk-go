package channel

//Status specifies status of the event
type Status int

const (
	//NEW if the event is new for the consumer
	NEW Status = iota
	// SUCCEED if the event is posted successfully
	SUCCESS
	//DELETE if the event is to delete
	DELETE
	//FAILED if the event  failed to post
	FAILED
)

//String represent of status enum
func (s Status) String() string {
	return [...]string{"NEW", "SUCCESS", "DELETE", "FAILED"}[s]
}

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

// String represent of Type enum
func (t Type) String() string {
	return [...]string{"LISTENER", "SENDER", "EVENT"}[t]
}
