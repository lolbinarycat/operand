package operand

import (
	"strconv"
)

type intValue int

func newIntValue(p *int, v int) *intValue {
	*p = v
	return (*intValue)(p)
}

func (v *intValue) Set(s string) error {
	num, err := strconv.Atoi(s)
	if err != nil {return err}
	*v = intValue(num)
	return nil
}

func (v intValue) String() string {
	return strconv.Itoa(int(v))
}

func (o *OpSet) IntVar(p *int, name string, value int, usage string) {
	o.Var(newIntValue(p,value),name,usage)
}

func (o *OpSet) Int(name string, value int, usage string) *int {
	p := new(int)
	o.IntVar(p,name,value,usage)
	return p
}
