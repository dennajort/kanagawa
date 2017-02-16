package bencode

import (
	"bufio"
	"fmt"
	"io"
)

func decodeIntegerLimit(r *bufio.Reader, e byte) (int64, error) {
	i := int64(0)
	for {
		c, err := r.ReadByte()
		switch {
		case err != nil, c == e:
			return i, err
		case c >= '0' && c <= '9':
			i = i*10 + int64(c-'0')
		default:
			return i, fmt.Errorf("Invalid character %c", c)
		}
	}
}

func decodeBytes(r *bufio.Reader) ([]byte, error) {
	le, err := decodeIntegerLimit(r, ':')
	if err != nil || le <= 0 {
		return nil, err
	}
	chunk := make([]byte, le)
	_, err = io.ReadFull(r, chunk)
	if err != nil {
		return nil, err
	}
	return chunk, nil
}

func decodeString(r *bufio.Reader) (string, error) {
	chunk, err := decodeBytes(r)
	return string(chunk), err
}

func decodeList(r *bufio.Reader) ([]interface{}, error) {
	l := make([]interface{}, 0)
	for {
		c, err := r.ReadByte()
		if err != nil || c == 'e' {
			return l, err
		}
		r.UnreadByte()
		elem, err := decodeAny(r)
		if err != nil {
			return l, err
		}
		l = append(l, elem)
	}
}

func decodeDict(r *bufio.Reader) (map[string]interface{}, error) {
	d := make(map[string]interface{})
	for {
		c, err := r.ReadByte()
		if err != nil || c == 'e' {
			return d, err
		}
		r.UnreadByte()
		key, err := decodeString(r)
		if err != nil {
			return d, err
		}
		value, err := decodeAny(r)
		if err != nil {
			return d, err
		}
		d[key] = value
	}
}

func decodeInteger(r *bufio.Reader) (int64, error) {
	sign := int64(1)
	c, err := r.ReadByte()
	switch {
	case err != nil:
		return 0, err
	case c == '-':
		sign = -1
	default:
		r.UnreadByte()
	}
	i, err := decodeIntegerLimit(r, 'e')
	return i * sign, err
}

func decodeAny(r *bufio.Reader) (interface{}, error) {
	c, err := r.ReadByte()
	if err != nil {
		return nil, err
	}
	switch c {
	case 'i':
		return decodeInteger(r)
	case 'l':
		return decodeList(r)
	case 'd':
		return decodeDict(r)
	default:
		r.UnreadByte()
		return decodeBytes(r)
	}
}

func Decode(r io.Reader) (interface{}, error) {
	bf := bufio.NewReader(r)
	return decodeAny(bf)
}
