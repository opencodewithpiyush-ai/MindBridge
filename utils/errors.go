package utils

type ErrInvalid struct {
	Field   string
	Message string
}

func (e *ErrInvalid) Error() string {
	return e.Message
}

// NewErrInvalid creates a new validation error.
func NewErrInvalid(field, message string) *ErrInvalid {
	return &ErrInvalid{Field: field, Message: message}
}
