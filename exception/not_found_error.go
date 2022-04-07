package exception

//NotFoundError to implement error interface
type NotFoundError struct {
	Error string
}

func NewNotFoundError(e string) NotFoundError {
	return  NotFoundError{Error: e}
}