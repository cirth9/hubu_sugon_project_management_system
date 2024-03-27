package errs

import "fmt"

type ErrorCode int

type BError struct {
	Code ErrorCode
	Msg  string
}

func (e *BError) Error() string {
	return fmt.Sprintf("code:%v,msg:%s", e.Code, e.Msg)
}

func NewError(code ErrorCode, msg string) *BError {
	return &BError{
		Code: code,
		Msg:  msg,
	}
}
