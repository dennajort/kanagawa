package bencode

import (
	"fmt"
	"io"
	"reflect"
	"runtime"
	"strconv"
)

type encoder struct {
	w interface {
		io.ByteWriter
		io.Writer
	}
}

func (e *encoder) writeByte(c byte) {
	if err := e.w.WriteByte(c); err != nil {
		panic(err)
	}
}

func (e *encoder) write(b []byte) {
	if _, err := e.w.Write(b); err != nil {
		panic(err)
	}
}

// marshalStructGetKey return the key for the bencode dict
// of the corresponding StructField
func marshalStructGetKey(field reflect.StructField) string {
	if field.PkgPath != "" {
		return ""
	}
	tag, exists := field.Tag.Lookup("benc")
	if !exists {
		return field.Name
	}
	switch tag {
	case "", "-":
		return ""
	default:
		return tag
	}
}

// marshalStruct create a bencode dict from a struct
func (e *encoder) marshalStruct(val reflect.Value) {
	e.writeByte('d')
	st := val.Type()
	for i := 0; i < st.NumField(); i++ {
		key := marshalStructGetKey(st.Field(i))
		if key == "" {
			continue
		}
		e.marshalString(key)
		e.marshalAny(val.Field(i))
	}
	e.writeByte('e')
}

// marshalMap create a bencode dict from a map[string]interface{}
func (e *encoder) marshalMap(val reflect.Value) {
	e.writeByte('d')
	for _, key := range val.MapKeys() {
		e.marshalString(key.String())
		e.marshalAny(val.MapIndex(key))
	}
	e.writeByte('e')
}

// marshalArrayOrSlice create a bencode list from a slice or array
func (e *encoder) marshalArrayOrSlice(val reflect.Value) {
	e.writeByte('l')
	for i := 0; i < val.Len(); i++ {
		e.marshalAny(val.Index(i))
	}
	e.writeByte('e')
}

func (e *encoder) marshalBytes(b []byte) {
	e.marshalInt(int64(len(b)))
	e.writeByte(':')
	e.write(b)
}

func (e *encoder) marshalString(s string) {
	e.marshalInt(int64(len(s)))
	e.writeByte(':')
	if _, err := io.WriteString(e.w, s); err != nil {
		panic(err)
	}
}

func (e *encoder) marshalInt(n int64) {
	e.write(strconv.AppendInt(nil, n, 10))
}

func (e *encoder) marshalIntOrUint(val reflect.Value) {
	e.writeByte('i')
	switch val.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		e.marshalInt(val.Int())
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		e.write(strconv.AppendUint(nil, val.Uint(), 10))
	default:
		panic(fmt.Errorf("Value is %s not Int or Uint", val.Kind()))
	}
	e.writeByte('e')
}

func (e *encoder) marshalAny(val reflect.Value) {
	switch k := val.Kind(); k {
	case reflect.String:
		e.marshalString(val.String())
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		e.marshalIntOrUint(val)
	case reflect.Slice:
		switch val.Type().Elem().Kind() {
		case reflect.Uint8:
			e.marshalBytes(val.Bytes())
		default:
			e.marshalArrayOrSlice(val)
		}
	case reflect.Array:
		e.marshalArrayOrSlice(val)
	case reflect.Map:
		switch val.Type().Key().Kind() {
		case reflect.String:
			e.marshalMap(val)
		default:
			panic(fmt.Errorf("Cannot reflect map key %s", val.Type().Key().Kind()))
		}
	case reflect.Interface, reflect.Ptr:
		e.marshalAny(val.Elem())
	case reflect.Struct:
		e.marshalStruct(val)
	default:
		panic(fmt.Errorf("Cannot reflect value %s", k))
	}
}

func (e *encoder) marshal(data interface{}) (err error) {
	defer func() {
		if e := recover(); e != nil {
			if _, ok := e.(runtime.Error); ok {
				panic(e)
			}
			err = e.(error)
		}
	}()

	e.marshalAny(reflect.ValueOf(data))
	return nil
}
