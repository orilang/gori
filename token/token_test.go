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

	t.Run("builtin_types", func(t *testing.T) {
		tests := []struct {
			input    Kind
			expected bool
		}{
			{
				input:    KWInt,
				expected: true,
			},
			{
				input:    KWFunc,
				expected: false,
			},
		}

		for _, tc := range tests {
			assert.Equal(tc.expected, IsBuiltinType(tc.input))
		}
	})

	t.Run("is_prefix", func(t *testing.T) {
		tests := []struct {
			input    Kind
			expected bool
		}{
			{
				input:    Plus,
				expected: true,
			},
			{
				input:    KWFunc,
				expected: false,
			},
		}

		for _, tc := range tests {
			assert.Equal(tc.expected, IsPrefix(tc.input))
		}
	})

	t.Run("is_infix", func(t *testing.T) {
		tests := []struct {
			input    Kind
			expected bool
		}{
			{
				input:    Lt,
				expected: true,
			},
			{
				input:    KWFunc,
				expected: false,
			},
		}

		for _, tc := range tests {
			assert.Equal(tc.expected, IsInfix(tc.input))
		}
	})

	t.Run("is_postfix", func(t *testing.T) {
		tests := []struct {
			input    Kind
			expected bool
		}{
			{
				input:    Dot,
				expected: true,
			},
			{
				input:    KWFunc,
				expected: false,
			},
		}

		for _, tc := range tests {
			assert.Equal(tc.expected, IsPostfix(tc.input))
		}
	})

	t.Run("is_comparison", func(t *testing.T) {
		tests := []struct {
			input    Kind
			expected bool
		}{
			{
				input:    Eq,
				expected: true,
			},
			{
				input:    KWFunc,
				expected: false,
			},
		}

		for _, tc := range tests {
			assert.Equal(tc.expected, IsComparison(tc.input))
		}
	})

	t.Run("is_chaining_comparison", func(t *testing.T) {
		tests := []struct {
			input    Kind
			expected bool
		}{
			{
				input:    Eq,
				expected: true,
			},
			{
				input:    KWFunc,
				expected: false,
			},
		}

		for _, tc := range tests {
			assert.Equal(tc.expected, IsChainingComparison(tc.input))
		}
	})

	t.Run("is_assignment", func(t *testing.T) {
		tests := []struct {
			input    Kind
			expected bool
		}{
			{
				input:    Assign,
				expected: true,
			},
			{
				input:    KWFunc,
				expected: false,
			},
		}

		for _, tc := range tests {
			assert.Equal(tc.expected, IsAssignment(tc.input))
		}
	})

	t.Run("is_range_for_assignment", func(t *testing.T) {
		tests := []struct {
			input    Kind
			expected bool
		}{
			{
				input:    Assign,
				expected: true,
			},
			{
				input:    KWFunc,
				expected: false,
			},
		}

		for _, tc := range tests {
			assert.Equal(tc.expected, IsRangeForAssignment(tc.input))
		}
	})

	t.Run("is_inc_dec", func(t *testing.T) {
		tests := []struct {
			input    Kind
			expected bool
		}{
			{
				input:    PPlus,
				expected: true,
			},
			{
				input:    KWFunc,
				expected: false,
			},
		}

		for _, tc := range tests {
			assert.Equal(tc.expected, IsIncDec(tc.input))
		}
	})
}
