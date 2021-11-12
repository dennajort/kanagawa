package bencode

import (
	"io/ioutil"
	"os"
	"testing"
)

func BenchmarkDecode(b *testing.B) {
	file, err := os.Open("../2dc18f47afee0307e138dab3015ee7e5154766f6.torrent")
	if err != nil {
		b.Fail()
	}
	buff, err := ioutil.ReadAll(file)
	if err != nil {
		b.Fail()
	}
	var data map[string]interface{}
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		err = Unmarshal(buff, &data)
		if err != nil {
			b.Fail()
		}
	}
}

func TestDecodeDict(t *testing.T) {
	t.Run("Empty dict", func(t *testing.T) {
		var d struct{}
		if err := UnmarshalString("de", &d); err != nil {
			t.Error(err)
		}
	})

	t.Run("Missing end", func(t *testing.T) {
		var d struct{}
		if err := UnmarshalString("d", &d); err == nil {
			t.Fail()
		}
	})

	t.Run("Dict of {foo:42}", func(t *testing.T) {
		d := struct {
			Foo int64
		}{}
		if err := UnmarshalString("d3:Fooi42ee", &d); err != nil {
			t.Error(err)
		}
		if d.Foo != 42 {
			t.Error("Int is not 42")
		}
	})

	t.Run("Dict of {Foo:42,Bar:24}", func(t *testing.T) {
		d := struct {
			Foo int64
			Bar int64
		}{}
		if err := UnmarshalString("d3:Fooi42e3:Bari24ee", &d); err != nil {
			t.Error(err)
		}
		if d.Foo != 42 {
			t.Error("Int is not 42")
		}
		if d.Bar != 24 {
			t.Error("Int is not 24")
		}
	})

	t.Run("Dict of {Foo:42,Bar:42} with tags", func(t *testing.T) {
		d := struct {
			Foo int64 `benc:"Bar"`
			Bar int64 `benc:"Foo"`
		}{}
		if err := UnmarshalString("d3:Fooi42e3:Bari24ee", &d); err != nil {
			t.Error(err)
		}
		if d.Foo != 24 {
			t.Error("Int is not 42")
		}
		if d.Bar != 42 {
			t.Error("Int is not 24")
		}
	})

	t.Run("Dict with wrong key", func(t *testing.T) {
		var d struct{}
		if err := UnmarshalString("di42e3:fooe", &d); err == nil {
			t.Fail()
		}
	})

	t.Run("Dict with wrong value", func(t *testing.T) {
		var d struct{}
		if err := UnmarshalString("d3:foofe", &d); err == nil {
			t.Fail()
		}
	})
}

func TestDecodeList(t *testing.T) {
	t.Run("Empty list", func(t *testing.T) {
		var l []int64
		if err := UnmarshalString("le", &l); err != nil {
			t.Error(err)
		}
		if len(l) != 0 {
			t.Fail()
		}
	})

	t.Run("Missing end", func(t *testing.T) {
		var l []int64
		if err := UnmarshalString("l", &l); err == nil {
			t.Fail()
		}
	})

	t.Run("List of [42]", func(t *testing.T) {
		var l []int64
		if err := UnmarshalString("li42ee", &l); err != nil {
			t.Error(err)
		}
		if len(l) != 1 {
			t.Fail()
		}
		if l[0] != 42 {
			t.Fail()
		}
	})

	t.Run("List of [42,'foo']", func(t *testing.T) {
		var l []interface{}
		if err := UnmarshalString("li42e3:fooe", &l); err != nil {
			t.Error(err)
		}
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
		var l []int64
		if err := UnmarshalString("lfe", &l); err == nil {
			t.Fail()
		}
	})
}

func TestDecodeString(t *testing.T) {
	t.Run("Empty string", func(t *testing.T) {
		s := "foo"
		if err := UnmarshalString("0:", &s); err != nil {
			t.Error(err)
		} else if s != "" {
			t.Fail()
		}
	})

	t.Run("Foo string", func(t *testing.T) {
		s := ""
		if err := UnmarshalString("3:Foo", &s); err != nil {
			t.Error(err)
		} else if s != "Foo" {
			t.Fail()
		}
	})

	t.Run("Empty string without length", func(t *testing.T) {
		s := "foo"
		if err := UnmarshalString(":", &s); err != nil {
			t.Error(err)
		} else if s != "" {
			t.Fail()
		}
	})

	t.Run("Too short string", func(t *testing.T) {
		s := ""
		if err := UnmarshalString("42:foo", &s); err == nil {
			t.Fail()
		}
	})
}

func TestDecodeInteger(t *testing.T) {
	t.Run("Empty integer", func(t *testing.T) {
		i := int64(42)
		if err := UnmarshalString("ie", &i); err != nil {
			t.Error(err)
		} else if i != 0 {
			t.Fail()
		}
	})

	t.Run("Zero integer", func(t *testing.T) {
		i := int64(42)
		if err := UnmarshalString("i0e", &i); err != nil {
			t.Error(err)
		} else if i != 0 {
			t.Fail()
		}
	})

	t.Run("Negative integer", func(t *testing.T) {
		i := int64(0)
		if err := UnmarshalString("i-42e", &i); err != nil {
			t.Error(err)
		} else if i != -42 {
			t.Fail()
		}
	})

	t.Run("Positive integer", func(t *testing.T) {
		i := int64(0)
		if err := UnmarshalString("i42e", &i); err != nil {
			t.Error(err)
		} else if i != 42 {
			t.Fail()
		}
	})

	t.Run("Unsigned integer", func(t *testing.T) {
		i := uint64(0)
		if err := UnmarshalString("i42e", &i); err != nil {
			t.Error(err)
		} else if i != 42 {
			t.Fail()
		}
	})

	t.Run("All number", func(t *testing.T) {
		i := int64(0)
		if err := UnmarshalString("i1234567890e", &i); err != nil {
			t.Error(err)
		} else if i != 1234567890 {
			t.Fail()
		}
	})

	t.Run("Invalid target", func(t *testing.T) {
		s := ""
		if err := UnmarshalString("i42", &s); err == nil {
			t.Fail()
		}
	})

	t.Run("Missing end", func(t *testing.T) {
		i := 0
		if err := UnmarshalString("i42", &i); err == nil {
			t.Fail()
		}
	})

	t.Run("Missing number and end", func(t *testing.T) {
		i := 0
		if err := UnmarshalString("i", &i); err == nil {
			t.Fail()
		}
	})

	t.Run("Integer as Bool", func(t *testing.T) {
		b := false
		if err := UnmarshalString("i1e", &b); err != nil {
			t.Error(err)
		} else if b != true {
			t.Fail()
		}
	})

	t.Run("Overflow an int", func(t *testing.T) {
		i := int8(0)
		if err := UnmarshalString("i1234567890e", &i); err == nil {
			t.Fail()
		}
	})

	t.Run("Overflow an uint", func(t *testing.T) {
		i := uint8(0)
		if err := UnmarshalString("i1234567890e", &i); err == nil {
			t.Fail()
		}
	})
}

func TestDecodeAny(t *testing.T) {
	t.Run("Empty", func(t *testing.T) {
		s := ""
		if err := UnmarshalString("", &s); err == nil {
			t.Fail()
		}
	})
}
