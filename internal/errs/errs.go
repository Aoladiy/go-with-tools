package errs

type Code int

const (
	InternalErrCode Code = iota
	UnauthorizedErrCode
	BadRequestErrCode
	UnprocessableEntityErrCode
	ConflictErrCode
	NotFoundErrCode
)

type AppError struct {
	Code Code
	Err  error
}

func (e *AppError) Error() string {
	if e.Code == InternalErrCode {
		return "internal server error"
	}
	return e.Err.Error()
}

func (e *AppError) Unwrap() error {
	return e.Err
}

func Internal(err error) *AppError {
	return &AppError{
		Code: InternalErrCode,
		Err:  err,
	}
}

func Unauthorized(err error) *AppError {
	return &AppError{
		Code: UnauthorizedErrCode,
		Err:  err,
	}
}

func BadRequest(err error) *AppError {
	return &AppError{
		Code: BadRequestErrCode,
		Err:  err,
	}
}

func UnprocessableEntity(err error) *AppError {
	return &AppError{
		Code: UnprocessableEntityErrCode,
		Err:  err,
	}
}

func Conflict(err error) *AppError {
	return &AppError{
		Code: ConflictErrCode,
		Err:  err,
	}
}

func NotFound(err error) *AppError {
	return &AppError{
		Code: NotFoundErrCode,
		Err:  err,
	}
}
