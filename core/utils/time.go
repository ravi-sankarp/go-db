package utils

import (
	"bytes"
	"time"
)

func TimeToBinary(t time.Time, buf *bytes.Buffer) error {
	if binaryData, err := t.UTC().MarshalBinary(); err != nil {
		return err
	} else {
		Serialize(binaryData, buf)
		return nil
	}
}

func BinaryToTime(binaryData []byte) (time.Time, error) {
	var deserializedTime time.Time
	err := deserializedTime.UnmarshalBinary(binaryData)
	return deserializedTime, err
}
