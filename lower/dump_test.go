package lower

import (
	"testing"

	"github.com/orilang/gori/ir"
)

func TestLower_dump(t *testing.T) {
	t.Run("nil", func(t *testing.T) {
		dump(nil)
	})

	t.Run("unhandled", func(t *testing.T) {
		dump(ir.Value("t"))
	})
}
