package bencode

import (
	"bufio"
	"bytes"
	"io"
	"strings"
)

func Decode(r io.Reader, v interface{}) error {
	d := &decoder{bufio.NewReader(r)}
	return d.unmarshal(v)
}

func Unmarshal(buff []byte, v interface{}) error {
	d := &decoder{bytes.NewReader(buff)}
	return d.unmarshal(v)
}

func UnmarshalString(buff string, v interface{}) error {
	d := &decoder{strings.NewReader(buff)}
	return d.unmarshal(v)
}
