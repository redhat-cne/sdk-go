package event

import (
	"strings"
)

//ValidationError ...
type ValidationError map[string]error

//Error ...
func (e ValidationError) Error() string {
	b := strings.Builder{}
	for k, v := range e {
		b.WriteString(k)
		b.WriteString(": ")
		b.WriteString(v.Error())
		b.WriteRune('\n')
	}
	return b.String()
}

// Validate performs a spec based validation on this event.
func (e Event) Validate() error {

	errs := map[string]error{}
	if e.FieldErrors != nil {
		for k, v := range e.FieldErrors {
			errs[k] = v
		}
	}

	if len(errs) > 0 {
		return ValidationError(errs)
	}
	return nil
}