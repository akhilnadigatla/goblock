package main

import (
	"bytes"
	"encoding/binary"
	"log"
)

// Convert an integer to its hexadecimal bytes
func intToHex(num int64) []byte {
	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.BigEndian, num)
	
	if err != nil {
		log.Panic(err)
	}
	
	return buf.Bytes()
}
