package utils

import (
	"bytes"
	"fmt"
	dbErrors "go-db/core/errors"
	"strconv"
	"strings"
	"time"
)

// Convert time.Time to a single string of fixed size
func TimeToString(t time.Time) (string, error) {

	if !t.Equal(t.UTC()) {
		return "", dbErrors.NewDbError("Timestamp should be in UTC")
	}

	utcTime := t.UTC()

	// Get the number of seconds and nanoseconds since Unix epoch (1 Jan 1970)
	secs := utcTime.Unix()
	nanos := utcTime.Nanosecond()

	// Convert to microseconds since Unix epoch
	microseconds := secs*1000000 + int64(nanos)/1000

	// Get the time zone offset in minutes
	_, offset := t.Zone()

	// Convert the microseconds timestamp and offset into strings
	timestampStr := strconv.FormatInt(microseconds, 10)
	offsetStr := strconv.Itoa(offset / 60) // convert offset to minutes

	// Combine timestamp and offset into a single string
	// Format: <timestamp>:<offset>
	result := timestampStr + ":" + offsetStr

	return result, nil
}

// Convert the combined string back to time.Time
func StringToTime(data string) (time.Time, error) {
	// Split the string into timestamp and offset parts
	parts := strings.Split(data, ":")
	if len(parts) != 2 {
		return time.Time{}, fmt.Errorf("invalid string format")
	}

	// Parse the timestamp and offset
	microseconds, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil {
		return time.Time{}, fmt.Errorf("failed to parse timestamp: %v", err)
	}

	offsetMinutes, err := strconv.Atoi(parts[1])
	if err != nil {
		return time.Time{}, fmt.Errorf("failed to parse offset: %v", err)
	}

	// Convert microseconds back to time.Time
	secs := microseconds / 1000000
	nanos := int((microseconds % 1000000) * 1000)
	utcTime := time.Unix(secs, int64(nanos))

	// Apply the offset to get the original time
	location := time.FixedZone("UTC", offsetMinutes*60)
	return utcTime.In(location), nil
}

func TimeToBinary(t time.Time, buf *bytes.Buffer) error {
	if str, err := TimeToString(t); err != nil {
		return err
	} else {
		Serialize([]byte(str), buf)
		return nil
	}
}

func BinaryToTime(byteArr []byte) (time.Time, error) {
	var timeString string
	DeSerialize(byteArr, &timeString, nil)
	return StringToTime(timeString)
}
