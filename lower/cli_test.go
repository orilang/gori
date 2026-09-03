package lower

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLower_NewLowerCLI(t *testing.T) {
	assert := assert.New(t)

	f := &Files{
		Files: []string{"xyz.ori"},
	}
	assert.Error(f.StartLowering())
}
