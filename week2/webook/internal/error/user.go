package error

import "fmt"

const (
	IncorrectUserNameOrPassword = iota + 10001
)

type UserError struct {
	err     error
	errCode int
}

func (ue UserError) Error() string {
	if ue.err == nil {
		return ""
	}
	return fmt.Sprintf("errcode:%d,msg:%s", ue.errCode, ue.err.Error())
}

func WrapError(err error, code int) error {
	if err == nil {
		return nil
	}
	return &UserError{err: err, errCode: code}
}
