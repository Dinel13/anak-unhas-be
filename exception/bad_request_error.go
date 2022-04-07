package exception

//NotFoundError to implement error interface
type BadRequestError struct {
	Error string
}

func NewBadRequestError(e string) BadRequestError {
	return  BadRequestError{Error: e}
}