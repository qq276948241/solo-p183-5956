package service

import "fmt"

type BizError struct {
	Code int
	Msg  string
}

func (e *BizError) Error() string {
	return fmt.Sprintf("biz error code=%d msg=%s", e.Code, e.Msg)
}

func NewBizError(code int, msg string) *BizError {
	return &BizError{Code: code, Msg: msg}
}
