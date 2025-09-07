package errors

type HttpError struct {
	Code int
	Msg  string
}

func (e *HttpError) Error() string {
	return e.Msg
}

var (
	ErrStartLine             = &HttpError{Code: 400, Msg: "Invalid Start Line"}
	ErrPartsMissingStartLine = &HttpError{Code: 400, Msg: "Some Parts of start line are missing"}
	ErrHttpPartsMissing      = &HttpError{Code: 400, Msg: "Some Parts of Http version are missing or http version not supported"}
)
