package custom_error

type SyntaxError struct {
	Message  string
	Position int
}

func (e SyntaxError) Error() string {
	return e.Message
}
