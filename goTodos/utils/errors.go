package utils

// HTTPError handles status codes
type HTTPError struct {
	Code int
	Msg  string
}

func (e HTTPError) Error() string {
	return e.Msg
}
