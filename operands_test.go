package operand

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestOpSet(t *testing.T) {
	asrt := assert.New(t)
	oSet := NewOpSet("test",ContinueOnError)
	tStr := oSet.String("t","123","test operand")
	err := oSet.Parse([]string{"t=7"})
	asrt.Nil(err)
	asrt.Equal(*tStr,"7")
}
