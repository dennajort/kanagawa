package bencode

import (
	"strings"
	"testing"
)

func TestDecodeInteger(t *testing.T) {
	t.Run("EmptyInteger", func(t *testing.T) {
		if i, err := Decode(strings.NewReader("ie")); err != nil {
			t.Error(err)
		} else if i.(int64) != 0 {
			t.Fail()
		}
	})
}
