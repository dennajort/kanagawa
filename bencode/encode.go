package bencode

import (
	"bufio"
	"bytes"
	"io"
)

func Encode(r io.Writer, v interface{}) error {
	wr := bufio.NewWriter(r)
	e := &encoder{wr}
	if err := e.marshal(v); err != nil {
		return err
	}
	return wr.Flush()
}

func Marshal(v interface{}) ([]byte, error) {
	buff := bytes.NewBuffer(nil)
	e := &encoder{buff}
	err := e.marshal(v)
	return buff.Bytes(), err
}

func MarshalString(v interface{}) (string, error) {
	buff := bytes.NewBuffer(nil)
	e := &encoder{buff}
	err := e.marshal(v)
	return buff.String(), err
}
