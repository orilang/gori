package semantic

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSemantic_NewTypeChecker(t *testing.T) {
	assert := assert.New(t)

	f := &Files{
		Files: []string{"xyz.ori"},
	}
	assert.Error(f.StartTypeChecking())
}
