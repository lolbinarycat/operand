// operand is a package for making command line tools with a similar interface to that of `dd`.
package operand

import (
	"fmt"
	"io"
	"os"
	"strings"
)

type ErrorHandling uint8

const (
	ContinueOnError ErrorHandling = iota
	_
	_
	ReturnOnError
)

type Value interface {
	String() string
	Set(string) error
}

// OpSet mirrors flag.FlagSet
type OpSet struct {
	Usage func()

	name string
	parsed bool
	ops map[string]*Operand
	output io.Writer
	errH ErrorHandling
}

func NewOpSet(name string,errorHandling ErrorHandling) *OpSet {
	o := &OpSet{
		name: name,
		errH: errorHandling,
	}
	// o.Usage = o.defaultUsage
	return o
}

func (o *OpSet) Parse(args []string) error {
	pErrs := NewParseErrors()
	for i, arg := range args {
		parts := strings.Split(arg,"=")
		if len(parts) != 2 {
			pErrs.Add(fmt.Errorf("wrong number '='s in arg %d",i))
			if o.errH == ReturnOnError {goto Return}
		} else {
			// key, value
			k, v := parts[0], parts[1]
			op, exists := o.ops[k]
			if exists {
				err := op.Value.Set(v)
				pErrs.Add(err)
			} else {
				pErrs.Add(NewPErr(fmt.Errorf("unknown operand %s",k),i))
				if o.errH == ReturnOnError {goto Return}
			}
		}
	}
Return:
	// the Canon method will return nil if
	// pErrs has no errors.
	return pErrs.Canon()
}

func (o *OpSet) Output() io.Writer {
	if o.output == nil {
		return os.Stderr
	}
	return o.output
}


type Operand struct {
	Name string
	Usage string
	Value Value
	DefValue string
}

func (o *OpSet) Var(value Value, name string, usage string) {
	op := &Operand{name,usage,value,value.String()}

	_, alreadythere := o.ops[name]
        if alreadythere {
                var msg string
                if o.name == "" {
                        msg = fmt.Sprintf("flag redefined: %s", name)
                } else {
                        msg = fmt.Sprintf("%s flag redefined: %s", o.name, name)
                }
                fmt.Fprintln(o.Output(), msg)
                panic(msg) // Happens only if flags are declared with identical names
        }
        if o.ops == nil {
                o.ops = make(map[string]*Operand)
        }
        o.ops[name] = op
}

func (o *OpSet) StringVar(p *string, name string, value string, usage string) {
	o.Var(newStringValue(p, value), name, usage)
}

func (o *OpSet) String(name, value, usage string) *string {
	p := new(string)
	o.StringVar(p,name,value,usage)
	return p
}

type stringValue string

func newStringValue(p *string, v string) *stringValue {
	*p = v
	return (*stringValue)(p)
}

func (s stringValue) String() string {
	return string(s)
}

func (s *stringValue) Set(str string) error {
	*s = stringValue(str)
	return nil
}
