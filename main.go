package main

import (
	"log"
	"os"

	"gitlab.com/dennajort/neptune/bencode"
	"gitlab.com/dennajort/neptune/metadata"
)

func main() {
	meta := metadata.Metadata{}
	err := bencode.Decode(os.Stdin, &meta)
	if err != nil {
		log.Panic(err)
	}
	// log.Println(meta)
	devNull, err := os.OpenFile(os.DevNull, os.O_WRONLY, 0755)
	if err != nil {
		log.Panic(err)
	}
	err = bencode.Encode(devNull, meta)
	if err != nil {
		log.Panic(err)
	}
	devNull.Close()
}
