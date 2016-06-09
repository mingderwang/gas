package validator

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

var (
	testData = map[string]string{
		"string": "test",
		"No":     "",
	}

	testRole = map[string]string{
		"string": "NotEmpty",
	}
)

func TestValidator(t *testing.T) {
	as := assert.New(t)

	v := Validator{}
	res, _ := v.Validate(testData, testRole)

	as.Equal(true, res)
}
