package codes

import (
	"fmt"
	"io"

	"github.com/dtynn/winston/internal/pb"
)

type codeError struct {
	cause error
	code  pb.ResultCode
}

func (c *codeError) Code() pb.ResultCode {
	return c.code
}

func (c *codeError) Error() string {
	return fmt.Sprintf("code %d; %s", c.code, c.cause)
}

func (c *codeError) Cause() error {
	return c.cause
}

func (c *codeError) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		if s.Flag('+') {
			fmt.Fprintf(s, "code %d\n%+v\n", c.Code(), c.Cause())
			return
		}
		fallthrough
	case 's', 'q':
		io.WriteString(s, c.Error())
	}
}

// WithCode alias a code to the error
func WithCode(code pb.ResultCode, err error) error {
	if err == nil {
		return nil
	}

	return &codeError{
		cause: err,
		code:  code,
	}
}

// Code return the code aliased to the error if any
func Code(err error) pb.ResultCode {
	if err == nil {
		return pb.ResultCode_ResultCodeOK
	}

	if e, ok := err.(*codeError); ok {
		return e.code
	}

	return pb.ResultCode_ResultCodeUnknown
}

// Result turn err into rpc result
func Result(err error) *pb.Result {
	res := &pb.Result{
		Code: Code(err),
	}

	if res.Code == pb.ResultCode_ResultCodeUnknown {
		res.Message = err.Error()
	}

	return res
}
