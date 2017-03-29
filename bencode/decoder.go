package bencode

import (
	"fmt"
	"io"
	"reflect"
	"runtime"
)

// Unmarshal argument must be a non-nil value of some pointer type.
type UnmarshalInvalidArgError struct {
	Type reflect.Type
}

func (e *UnmarshalInvalidArgError) Error() string {
	if e.Type == nil {
		return "bencode: Unmarshal(nil)"
	}

	if e.Type.Kind() != reflect.Ptr {
		return fmt.Sprintf("bencode: Unmarshal(non-pointer %s)", e.Type)
	}
	return fmt.Sprintf("bencode: Unmarshal(nil %s)", e.Type)
}

// Unmarshaler spotted a value that was not appropriate for a given Go value.
type UnmarshalTypeError struct {
	Value string
	Type  reflect.Type
}

func (e *UnmarshalTypeError) Error() string {
	return "bencode: value (" + e.Value + ") is not appropriate for type: " +
		e.Type.String()
}

type decoder struct {
	r interface {
		io.ByteScanner
		io.Reader
	}
}

func (d *decoder) readByte() byte {
	c, err := d.r.ReadByte()
	if err != nil {
		panic(err)
	}
	return c
}

func (d *decoder) unreadByte() {
	if err := d.r.UnreadByte(); err != nil {
		panic(err)
	}
}

func (d *decoder) peekByte(c byte) bool {
	if d.readByte() == c {
		return true
	}
	d.unreadByte()
	return false
}

func (d *decoder) readBytes(le int) []byte {
	buff := make([]byte, le)
	if _, err := io.ReadFull(d.r, buff); err != nil {
		panic(err)
	}
	return buff
}

func (d *decoder) decodeUintLimit(e byte) uint64 {
	i := uint64(0)
	for {
		c := d.readByte()
		switch {
		case c == e:
			return i
		case c >= '0' && c <= '9':
			i = i*10 + uint64(c-'0')
		default:
			panic(fmt.Errorf("Invalid character %c", c))
		}
	}
}

func (d *decoder) decodeBytes() []byte {
	le := d.decodeUintLimit(':')
	return d.readBytes(int(le))
}

func (d *decoder) decodeBytesOrString(pv reflect.Value) {
	chunk := d.decodeBytes()
	switch {
	case pv.Kind() == reflect.String:
		pv.SetString(string(chunk))
	case pv.Kind() == reflect.Slice && pv.Type().Elem().Kind() == reflect.Uint8:
		pv.SetBytes(chunk)
	default:
		panic(&UnmarshalTypeError{
			Value: "string",
			Type:  pv.Type(),
		})
	}
}

func (d *decoder) decodeList(pv reflect.Value) {
	switch pv.Kind() {
	case reflect.Slice:
	default:
		panic(&UnmarshalTypeError{
			Value: "array",
			Type:  pv.Type(),
		})
	}

	childType := pv.Type().Elem()
	for i := 0; !d.peekByte('e'); i++ {
		if i >= pv.Len() {
			pv.Set(reflect.Append(pv, reflect.Zero(childType)))
		}
		d.decodeAny(pv.Index(i))
	}
}

func (d *decoder) decodeListValue() (l []interface{}) {
	for !d.peekByte('e') {
		l = append(l, d.decodeAnyValue())
	}
	return
}

func (d *decoder) decodeDictMap(pv reflect.Value) {
	mapType := pv.Type()
	childType := mapType.Elem()

	if pv.IsNil() {
		pv.Set(reflect.MakeMap(mapType))
	}

	for !d.peekByte('e') {
		var key string
		d.decodeBytesOrString(reflect.ValueOf(&key).Elem())
		val := reflect.New(childType).Elem()
		d.decodeAny(val)
		pv.SetMapIndex(reflect.ValueOf(key), val)
	}
}

func (d *decoder) decodeDictStruct(pv reflect.Value) {
	st := pv.Type()
	mapping := make(map[string]int)
	for i := 0; i < st.NumField(); i++ {
		field := st.Field(i)
		if field.PkgPath != "" {
			continue
		}
		tag, exists := field.Tag.Lookup("benc")
		if !exists {
			mapping[field.Name] = i
			continue
		}
		switch tag {
		case "", "-":
			break
		default:
			mapping[tag] = i
		}
	}
	for !d.peekByte('e') {
		var key string
		d.decodeBytesOrString(reflect.ValueOf(&key).Elem())
		idx, exists := mapping[key]
		if exists {
			f := pv.Field(idx)
			d.decodeAny(f)
		} else {
			d.decodeAnyValue()
		}
	}
}

func (d *decoder) decodeDict(pv reflect.Value) {
	k := pv.Kind()
	switch {
	case k == reflect.Struct:
		d.decodeDictStruct(pv)
	case k == reflect.Map && pv.Type().Key().Kind() == reflect.String:
		d.decodeDictMap(pv)
	default:
		panic(&UnmarshalTypeError{
			Value: "dict",
			Type:  pv.Type(),
		})
	}
}

func (d *decoder) decodeDictValue() map[string]interface{} {
	m := make(map[string]interface{})
	for !d.peekByte('e') {
		key := string(d.decodeBytes())
		m[key] = d.decodeAnyValue()
	}
	return m
}

func (d *decoder) decodeIntegerBool(pv reflect.Value) {
	pv.SetBool(d.decodeUintLimit('e') != 0)
}

func (d *decoder) decodeIntValue() int64 {
	val := int64(1)
	if d.peekByte('-') {
		val = -1
	}
	return val * int64(d.decodeUintLimit('e'))
}

func (d *decoder) decodeIntegerInt(pv reflect.Value) {
	val := d.decodeIntValue()
	if pv.OverflowInt(val) {
		panic(&UnmarshalTypeError{
			Value: fmt.Sprintf("integer %d", val),
			Type:  pv.Type(),
		})
	}
	pv.SetInt(val)
}

func (d *decoder) decodeIntegerUint(pv reflect.Value) {
	i := d.decodeUintLimit('e')
	if pv.OverflowUint(i) {
		panic(&UnmarshalTypeError{
			Value: fmt.Sprintf("integer %d", i),
			Type:  pv.Type(),
		})
	}
	pv.SetUint(i)
}

func (d *decoder) decodeInteger(pv reflect.Value) {
	switch pv.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		d.decodeIntegerInt(pv)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		d.decodeIntegerUint(pv)
	case reflect.Bool:
		d.decodeIntegerBool(pv)
	default:
		panic(&UnmarshalTypeError{
			Value: "integer",
			Type:  pv.Type(),
		})
	}
}

func (d *decoder) decodeAny(pv reflect.Value) {
	if pv.Kind() == reflect.Ptr {
		// if the pointer is nil, allocate a new element of the type it points to
		if pv.IsNil() {
			pv.Set(reflect.New(pv.Type().Elem()))
		}
		pv = pv.Elem()
	}

	if pv.Kind() == reflect.Interface {
		pv.Set(reflect.ValueOf(d.decodeAnyValue()))
		return
	}

	switch d.readByte() {
	case 'i':
		d.decodeInteger(pv)
	case 'l':
		d.decodeList(pv)
	case 'd':
		d.decodeDict(pv)
	default:
		d.unreadByte()
		d.decodeBytesOrString(pv)
	}
}

func (d *decoder) decodeAnyValue() interface{} {
	switch d.readByte() {
	case 'i':
		return d.decodeIntValue()
	case 'l':
		return d.decodeListValue()
	case 'd':
		return d.decodeDictValue()
	default:
		d.unreadByte()
		return d.decodeBytes()
	}
}

func (d *decoder) unmarshal(v interface{}) (err error) {
	defer func() {
		if e := recover(); e != nil {
			if _, ok := e.(runtime.Error); ok {
				panic(e)
			}
			err = e.(error)
		}
	}()

	pv := reflect.ValueOf(v)
	if pv.Kind() != reflect.Ptr || pv.IsNil() {
		return &UnmarshalInvalidArgError{reflect.TypeOf(v)}
	}
	d.decodeAny(pv.Elem())
	return nil
}
