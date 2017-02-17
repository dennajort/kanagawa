package bencode

import (
	"bytes"
	"strings"
	"testing"
)

func TestDecodeDict(t *testing.T) {
	t.Run("Empty dict", func(t *testing.T) {
		i, err := Decode(strings.NewReader("de"))
		if err != nil {
			t.Error(err)
		}
		d := i.(map[string]interface{})
		if len(d) != 0 {
			t.Fail()
		}
	})

	t.Run("Missing end", func(t *testing.T) {
		if _, err := Decode(strings.NewReader("d")); err == nil {
			t.Fail()
		}
	})

	t.Run("Dict of {foo:42}", func(t *testing.T) {
		i, err := Decode(strings.NewReader("d3:fooi42ee"))
		if err != nil {
			t.Error(err)
		}
		d := i.(map[string]interface{})
		if len(d) != 1 {
			t.Fail()
		}
		if d["foo"].(int64) != 42 {
			t.Error("Int is not 42")
		}
	})

	t.Run("Dict of {foo:42,bar:24}", func(t *testing.T) {
		i, err := Decode(strings.NewReader("d3:fooi42e3:bari24ee"))
		if err != nil {
			t.Error(err)
		}
		d := i.(map[string]interface{})
		if len(d) != 2 {
			t.Fail()
		}
		if d["foo"].(int64) != 42 {
			t.Error("Int is not 42")
		}
		if d["bar"].(int64) != 24 {
			t.Error("Int is not 24")
		}
	})

	t.Run("Dict with wrong key", func(t *testing.T) {
		if _, err := Decode(strings.NewReader("di42e3:fooe")); err == nil {
			t.Fail()
		}
	})

	t.Run("Dict with wrong value", func(t *testing.T) {
		if _, err := Decode(strings.NewReader("d3:foofe")); err == nil {
			t.Fail()
		}
	})
}

func TestDecodeList(t *testing.T) {
	t.Run("Empty list", func(t *testing.T) {
		i, err := Decode(strings.NewReader("le"))
		if err != nil {
			t.Error(err)
		}
		l := i.([]interface{})
		if len(l) != 0 {
			t.Fail()
		}
	})

	t.Run("Missing end", func(t *testing.T) {
		if _, err := Decode(strings.NewReader("l")); err == nil {
			t.Fail()
		}
	})

	t.Run("List of [42]", func(t *testing.T) {
		i, err := Decode(strings.NewReader("li42ee"))
		if err != nil {
			t.Error(err)
		}
		l := i.([]interface{})
		if len(l) != 1 {
			t.Fail()
		}
		if l[0].(int64) != 42 {
			t.Fail()
		}
	})

	t.Run("List of [42,'foo']", func(t *testing.T) {
		i, err := Decode(strings.NewReader("li42e3:fooe"))
		if err != nil {
			t.Error(err)
		}
		l := i.([]interface{})
		if len(l) != 2 {
			t.Fail()
		}
		if l[0].(int64) != 42 {
			t.Error("Int is not 42")
		}
		if string(l[1].([]byte)) != "foo" {
			t.Error("String is not foo")
		}
	})

	t.Run("List with wrong item", func(t *testing.T) {
		if _, err := Decode(strings.NewReader("lfe")); err == nil {
			t.Fail()
		}
	})
}

func TestDecodeString(t *testing.T) {
	t.Run("Empty string", func(t *testing.T) {
		if i, err := Decode(strings.NewReader("0:")); err != nil {
			t.Error(err)
		} else if !bytes.Equal(i.([]byte), []byte{}) {
			t.Fail()
		}
	})

	t.Run("Foo string", func(t *testing.T) {
		if i, err := Decode(strings.NewReader("3:Foo")); err != nil {
			t.Error(err)
		} else if string(i.([]byte)) != "Foo" {
			t.Fail()
		}
	})

	t.Run("Empty string without length", func(t *testing.T) {
		if i, err := Decode(strings.NewReader(":")); err != nil {
			t.Error(err)
		} else if string(i.([]byte)) != "" {
			t.Fail()
		}
	})

	t.Run("Too short string", func(t *testing.T) {
		if _, err := Decode(strings.NewReader("42:foo")); err == nil {
			t.Fail()
		}
	})
}

func TestDecodeInteger(t *testing.T) {
	t.Run("Empty integer", func(t *testing.T) {
		if i, err := Decode(strings.NewReader("ie")); err != nil {
			t.Error(err)
		} else if i.(int64) != 0 {
			t.Fail()
		}
	})

	t.Run("Zero integer", func(t *testing.T) {
		if i, err := Decode(strings.NewReader("i0e")); err != nil {
			t.Error(err)
		} else if i.(int64) != 0 {
			t.Fail()
		}
	})

	t.Run("Negative integer", func(t *testing.T) {
		if i, err := Decode(strings.NewReader("i-42e")); err != nil {
			t.Error(err)
		} else if i.(int64) != -42 {
			t.Fail()
		}
	})

	t.Run("Positive integer", func(t *testing.T) {
		if i, err := Decode(strings.NewReader("i42e")); err != nil {
			t.Error(err)
		} else if i.(int64) != 42 {
			t.Fail()
		}
	})

	t.Run("All number", func(t *testing.T) {
		if i, err := Decode(strings.NewReader("i1234567890e")); err != nil {
			t.Error(err)
		} else if i.(int64) != 1234567890 {
			t.Fail()
		}
	})

	t.Run("Missing end", func(t *testing.T) {
		if _, err := Decode(strings.NewReader("i42")); err == nil {
			t.Fail()
		}
	})

	t.Run("Missing number and end", func(t *testing.T) {
		if _, err := Decode(strings.NewReader("i")); err == nil {
			t.Fail()
		}
	})
}

func TestDecodeAny(t *testing.T) {
	t.Run("Empty", func(t *testing.T) {
		if _, err := Decode(strings.NewReader("")); err == nil {
			t.Fail()
		}
	})
}
