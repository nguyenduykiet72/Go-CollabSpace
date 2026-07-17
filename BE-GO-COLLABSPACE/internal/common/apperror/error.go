package apperror

type AppError struct {
	StatusCode int    `json:"-"`
	RootErr    error  `json:"-"`
	Message    string `json:"message"`
	Log        string `json:"log"`
	Key        string `json:"error_key"`
}

func (e *AppError) Error() string {
	return e.Message
}

func (e *AppError) Unwrap() error {
	return e.RootErr
}

func (e *AppError) RootError() error {
	if e.RootErr != nil {
		return e.RootErr
	}
	return e
}

func (e *AppError) WithRootErr(err error) *AppError {
	newErr := *e
	newErr.RootErr = err
	newErr.Log = err.Error()
	return &newErr
}

func (e *AppError) WithLog(log string) *AppError {
	newErr := *e
	newErr.Log = log
	return &newErr
}

func (e *AppError) WithMessage(msg string) *AppError {
	newErr := *e
	newErr.Message = msg
	return &newErr
}

func NewAppError(statusCode int, msg, key string) *AppError {
	return &AppError{
		StatusCode: statusCode,
		Message:    msg,
		Key:        key,
	}
}
