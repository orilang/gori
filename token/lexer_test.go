package token

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestToken(t *testing.T) {
	assert := assert.New(t)

	t.Run("token", func(t *testing.T) {
		tests := []struct {
			input    string
			expected Kind
		}{
			{
				input:    "true",
				expected: BoolLit,
			},
			{
				input:    "gori",
				expected: Ident,
			},
		}

		for _, tc := range tests {
			assert.Equal(tc.expected, LookupKeyword(tc.input))
		}
	})
}
