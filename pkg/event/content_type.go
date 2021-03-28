package event

const (
	// TextPlain ...
	TextPlain = "text/plain"
	// TextJSON ...
	TextJSON = "text/json"
	// ApplicationJSON ...
	ApplicationJSON = "application/json"
)

// StringOfApplicationJSON returns a string pointer to "application/json"
func StringOfApplicationJSON() *string {
	a := ApplicationJSON
	return &a
}

// StringOfTextPlain returns a string pointer to "text/plain"
func StringOfTextPlain() *string {
	a := TextPlain
	return &a
}
