package utils

import (
	"bytes"
)

func CheckIfString(value any) ([]byte, bool) {
	if str, ok := value.(string); ok == true {
		return []byte(str), ok
	} else {
		return nil, ok
	}
}

func StringToBinary(value string, buf *bytes.Buffer) {
	strBuf := []byte(value)
	// adding length header
	Serialize(int32(len(strBuf)), buf)
	Serialize(strBuf, buf)
}
