package utils

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

func Serialize(value any, buf *bytes.Buffer) {
	if err := binary.Write(buf, binary.LittleEndian, value); err != nil {
		fmt.Println(err)
		panic(fmt.Sprint("Serialization failed : ", err))
	}
}

func DeSerialize(buffer []byte, target any, reader *bytes.Reader) {
	if reader == nil {
		reader = bytes.NewReader(buffer)
	}

	if err := binary.Read(reader, binary.LittleEndian, target); err != nil {
		fmt.Println(err)
		panic(fmt.Sprint("Deserialization failed :", err))
	}
}
