package bencode

import (
	"bytes"
	"io/ioutil"
	"os"
	"testing"
)

func BenchmarkEncode(b *testing.B) {
	file, err := os.Open("../2dc18f47afee0307e138dab3015ee7e5154766f6.torrent")
	if err != nil {
		b.Fatal(err)
	}
	buff, err := ioutil.ReadAll(file)
	if err != nil {
		b.Fatal(err)
	}
	var data map[string]interface{}
	err = Unmarshal(buff, &data)
	if err != nil {
		b.Fatal(err)
	}
	devNull, err := os.OpenFile(os.DevNull, os.O_WRONLY, 0755)
	if err != nil {
		b.Fatal(err)
	}
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		err := Encode(devNull, data)
		if err != nil {
			b.Fail()
		}
	}
}

func TestMarshalNil(t *testing.T) {
	t.Run("Nil", func(t *testing.T) {
		if _, err := Marshal(nil); err == nil {
			t.Fail()
		}
	})
}

func TestMarshalArray(t *testing.T) {
	t.Run("Empty Slice", func(t *testing.T) {
		s, err := MarshalString([]int{})
		if err != nil {
			t.Error(err)
		}
		if s != "le" {
			t.Fail()
		}
	})

	t.Run("One Slice", func(t *testing.T) {
		s, err := MarshalString([]int{42})
		if err != nil {
			t.Error(err)
		}
		if s != "li42ee" {
			t.Fail()
		}
	})

	t.Run("One Array", func(t *testing.T) {
		s, err := MarshalString([...]int{42})
		if err != nil {
			t.Error(err)
		}
		if s != "li42ee" {
			t.Fail()
		}
	})

	t.Run("Nil array", func(t *testing.T) {
		b := bytes.NewBuffer(nil)
		if err := Decode(b, []interface{}{nil}); err == nil {
			t.Fail()
		}
	})
}

func TestMarshalString(t *testing.T) {
	t.Run("String Foo", func(t *testing.T) {
		s, err := MarshalString("Foo")
		if err != nil {
			t.Error(err)
		}
		if s != "3:Foo" {
			t.Fail()
		}
	})

	t.Run("Slice Foo", func(t *testing.T) {
		s, err := MarshalString([]byte("Foo"))
		if err != nil {
			t.Error(err)
		}
		if s != "3:Foo" {
			t.Fail()
		}
	})
}

func TestMarshalInt(t *testing.T) {
	t.Run("Int 42", func(t *testing.T) {
		s, err := MarshalString(int(42))
		if err != nil {
			t.Error(err)
		}
		if s != "i42e" {
			t.Fail()
		}
	})

	t.Run("Int -42", func(t *testing.T) {
		s, err := MarshalString(int(-42))
		if err != nil {
			t.Error(err)
		}
		if s != "i-42e" {
			t.Fail()
		}
	})
}

func TestMarshalUint(t *testing.T) {
	t.Run("Uint 42", func(t *testing.T) {
		n := uint(42)
		s, err := MarshalString(&n)
		if err != nil {
			t.Error(err)
		}
		if s != "i42e" {
			t.Fail()
		}
	})
}
