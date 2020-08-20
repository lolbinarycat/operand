package operand

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type ParseError struct {
	Err error
	Arg int
}

func (e *ParseError) Error() string {
	return fmt.Sprintf("error in arg %d: %s",e.Arg+1,e.Err)
}

func (e *ParseError) Unwrap() error {
	return e.Err
}

func NewPErr(err error,arg int) *ParseError {
	return &ParseError{
		Err: err,
		Arg: arg,
	}
}

func NewStrPErr(str string,arg int) *ParseError {
	return &ParseError{
		Err: errors.New(str),
		Arg: arg,
	}
}


type ParseErrors struct {
	Errors []error
}

func NewParseErrors() *ParseErrors {
	return &ParseErrors{
		Errors: make([]error,0,4),
	}
}

func (e *ParseErrors) Error() string {
	var bldr strings.Builder
	bldr.WriteString("errors:")
	for i, err := range e.Errors {
		bldr.WriteRune('\n')
		bldr.WriteString(strconv.Itoa(i))
		bldr.WriteRune(':')
		bldr.WriteRune(' ')
		bldr.WriteString(err.Error())
	}
	return bldr.String()
}

// Cannon returns nil if e contains no errors
func (e *ParseErrors) Canon() *ParseErrors {
	if len(e.Errors) == 0 {
		return nil
	}
	return e
}

// Add adds err to e if err isn't nil.
func (e *ParseErrors) Add(err error) {
	if err != nil {
		e.Errors = append(e.Errors,err) 
	}
}

// Is returns true if errors.Is(err, target)
// returns true for every element of e.
func (e *ParseErrors) Is(target error) bool {
	if e == target {return true}
	if len(e.Errors) == 0 {return target == nil}
	
	for _, err := range e.Errors {
		if errors.Is(err, target) == false {
			return false
		}
	}
	return true
}
