package bencode

import (
	"bytes"
	"testing"
)

func TestMarshalInt(t *testing.T) {
	t.Run("Int 42", func(t *testing.T) {
		b := bytes.NewBuffer(nil)
		if err := Marshal(b, int(42)); err != nil {
			t.Error(err)
		}
		if b.String() != "i42e" {
			t.Fail()
		}
	})

	t.Run("Int -42", func(t *testing.T) {
		b := bytes.NewBuffer(nil)
		if err := Marshal(b, int(-42)); err != nil {
			t.Error(err)
		}
		if b.String() != "i-42e" {
			t.Fail()
		}
	})
}

func TestMarshalUint(t *testing.T) {
	t.Run("Uint 42", func(t *testing.T) {
		b := bytes.NewBuffer(nil)
		if err := Marshal(b, uint(42)); err != nil {
			t.Error(err)
		}
		if b.String() != "i42e" {
			t.Fail()
		}
	})
}
