package exception

type Error interface {
	error
	Code() int
	Msg() string
}

type NotFoundHttpError struct{}

func (e NotFoundHttpError) Error() string {
	return "Http not found"
}

type MethodNotAllowedHttpError struct{}

func (e MethodNotAllowedHttpError) Error() string {
	return "Method not allowed"
}
