package bencode

import (
	"bufio"
	"fmt"
	"io"
	"reflect"
	"strconv"
)

func marshalMap(w *bufio.Writer, val reflect.Value) (err error) {
	if err = w.WriteByte('d'); err != nil {
		return
	}
	for _, key := range val.MapKeys() {
		if err = marshalStringOrBytes(w, key); err != nil {
			return
		}
		if err = marshalAny(w, val.MapIndex(key)); err != nil {
			return
		}
	}
	return w.WriteByte('e')
}

func marshalArrayOrSlice(w *bufio.Writer, val reflect.Value) (err error) {
	if err = w.WriteByte('l'); err != nil {
		return
	}
	for i := 0; i < val.Len(); i++ {
		if err = marshalAny(w, val.Index(i)); err != nil {
			return
		}
	}
	return w.WriteByte('e')
}

func marshalStringOrBytes(w *bufio.Writer, val reflect.Value) (err error) {
	if _, err = w.Write(strconv.AppendInt(nil, int64(val.Len()), 10)); err != nil {
		return
	}
	if err = w.WriteByte(':'); err != nil {
		return
	}
	switch val.Kind() {
	case reflect.String:
		_, err = w.WriteString(val.String())
	default:
		_, err = w.Write(val.Bytes())
	}
	return
}

func marshalIntOrUint(w *bufio.Writer, val reflect.Value) (err error) {
	if err = w.WriteByte('i'); err != nil {
		return
	}
	switch k := val.Kind(); k {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		_, err = w.Write(strconv.AppendInt(nil, val.Int(), 10))
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		_, err = w.Write(strconv.AppendUint(nil, val.Uint(), 10))
	default:
		return fmt.Errorf("Value is %s not Int or Uint", k)
	}
	if err != nil {
		return
	}
	return w.WriteByte('e')
}

func marshalAny(w *bufio.Writer, val reflect.Value) error {
	switch k := val.Kind(); k {
	case reflect.String:
		return marshalStringOrBytes(w, val)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return marshalIntOrUint(w, val)
	case reflect.Slice:
		switch val.Type().Elem().Kind() {
		case reflect.Uint8:
			return marshalStringOrBytes(w, val)
		default:
			return marshalArrayOrSlice(w, val)
		}
	case reflect.Array:
		return marshalArrayOrSlice(w, val)
	case reflect.Map:
		return marshalMap(w, val)
	case reflect.Interface:
		return marshalAny(w, val.Elem())
	default:
		return fmt.Errorf("Cannot reflect value %s", k)
	}
}

func Marshal(w io.Writer, data interface{}) (err error) {
	bw := bufio.NewWriter(w)
	if err = marshalAny(bw, reflect.ValueOf(data)); err != nil {
		return
	}
	return bw.Flush()
}
