package main

// InputError ...
type InputError struct {
	Err string
}

func (e *InputError) Error() string {
	return e.Err
}

// NewInputError ...
func NewInputError(err string) error {
	return &InputError{
		Err: err,
	}
}

// NotFoundError ...
type NotFoundError struct {
	Err string
}

func (e *NotFoundError) Error() string {
	return e.Err
}

// NewNotFoundError ...
func NewNotFoundError(err string) error {
	return &NotFoundError{
		Err: err,
	}
}
