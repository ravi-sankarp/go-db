package internal

import (
	"bytes"
	"fmt"
	"go-db/core/constants"
	dbErrors "go-db/core/errors"
	"go-db/core/utils"
	"reflect"
	"time"
)

var SAMPLE_TUPLE_SCHEMA = Tuple{
	attribute{
		name:     "name",
		value:    reflect.ValueOf(constants.STRING),
		dataType: reflect.String,
	},
	attribute{
		name:     "email",
		value:    reflect.ValueOf(constants.STRING),
		dataType: reflect.String,
	},
	attribute{
		name:     "age",
		value:    reflect.ValueOf(constants.INTEGER),
		dataType: reflect.String,
	},
	attribute{
		name:     "deleted",
		value:    reflect.ValueOf(constants.BOOLEAN),
		dataType: reflect.String,
	},
	attribute{
		name:     "created_on",
		value:    reflect.ValueOf(constants.TIMESTAMP),
		dataType: reflect.String,
	},
}

type SAMPLE_TABLE_STRUCT struct {
	Name       string
	Email      string
	Age        int
	Deleted    bool
	Created_on time.Time
}

var SAMPLE_TUPLE_INSERT = Tuple{
	attribute{
		name:     "name",
		value:    reflect.ValueOf("test"),
		dataType: reflect.String,
	},
	attribute{
		name:     "email",
		value:    reflect.ValueOf("test@gmail.com"),
		dataType: reflect.String,
	},
	attribute{
		name:     "age",
		value:    reflect.ValueOf(23),
		dataType: reflect.Int,
	},
	attribute{
		name:     "deleted",
		value:    reflect.ValueOf(false),
		dataType: reflect.Bool,
	},
	attribute{
		name:     "created_on",
		value:    reflect.ValueOf(time.Now().UTC()),
		dataType: reflect.Struct,
	},
}

func parseDataToBinaryTuple(tuple Tuple) (*bytes.Buffer, error) {
	var (
		dataBuf = new(bytes.Buffer)
	)
	for _, attribute := range tuple {
		switch attribute.dataType {
		case reflect.String:
			utils.StringToBinary(attribute.value.Interface().(string), dataBuf)
		case reflect.Int64:
		case reflect.Int:
			utils.Serialize(attribute.value.Int(), dataBuf)
		case reflect.Bool:
			utils.Serialize(attribute.value.Interface().(bool), dataBuf)
		case reflect.Struct:
			if reflect.TypeOf(attribute.value.Interface()) == reflect.TypeOf(time.Time{}) {
				value := attribute.value.Interface().(time.Time)
				utils.TimeToBinary(value, dataBuf)
			} else {
				fmt.Println(attribute.dataType)
				fmt.Println(reflect.TypeOf(attribute.dataType.String()))
				return dataBuf, dbErrors.NewDbError("Invalid Type for " + attribute.name)
			}
		default:
			return dataBuf, dbErrors.NewDbError("Invalid Type for " + attribute.name)
		}
	}
	return dataBuf, nil
}

func parseTupleFromPageBuffer(buf []byte) (Tuple, error) {
	var bytesRead = 0
	var tuple = make(Tuple, len(SAMPLE_TUPLE_SCHEMA))
	for i, attr := range SAMPLE_TUPLE_SCHEMA {
		switch attr.value.String() {
		case constants.STRING:
			// get string length header and parse the data
			var stringSize stringHeader
			utils.DeSerialize(buf[bytesRead:bytesRead+int(constants.STRING_HEADER_SIZE)], &stringSize, nil)
			var strBuf = make([]byte, stringSize.Length)
			utils.DeSerialize(buf[bytesRead+int(constants.STRING_HEADER_SIZE):bytesRead+int(constants.STRING_HEADER_SIZE)+int(stringSize.Length)], &strBuf, nil)

			tuple[i] = attribute{
				name:     attr.name,
				value:    reflect.ValueOf(string(strBuf)),
				dataType: reflect.String,
			}
			bytesRead += int(constants.STRING_HEADER_SIZE) + int(stringSize.Length)

		case constants.INTEGER:
			var intValue int64
			utils.DeSerialize(buf[bytesRead:bytesRead+constants.INT_SIZE], &intValue, nil)

			tuple[i] = attribute{
				name:     attr.name,
				value:    reflect.ValueOf(intValue),
				dataType: reflect.Int64,
			}
			bytesRead += constants.INT_SIZE

		case constants.BOOLEAN:
			var boolValue bool
			utils.DeSerialize(buf[bytesRead:bytesRead+constants.BOOL_SIZE], &boolValue, nil)

			tuple[i] = attribute{
				name:     attr.name,
				value:    reflect.ValueOf(boolValue),
				dataType: reflect.Bool,
			}
			bytesRead += int(constants.BOOL_SIZE)

		case constants.TIMESTAMP:

			timeValue, err := utils.BinaryToTime(buf[bytesRead : bytesRead+constants.TIMESTAMP_SIZE])
			if err != nil {
				return nil, dbErrors.NewDbError("Error while parsing attribute " + attr.name + " error: " + err.Error())
			}
			tuple[i] = attribute{
				name:     attr.name,
				value:    reflect.ValueOf(timeValue),
				dataType: reflect.Struct,
			}
			bytesRead += int(constants.BOOL_SIZE)

		default:
			return nil, dbErrors.NewDbError("Invalid Type while parsing page for attribute " + attr.name)
		}
	}
	return tuple, nil
}

func parseTuple(reader *bytes.Reader, byteOffset uint16) (Tuple, error) {
	err := seekBufferReader(reader, int64(byteOffset), 0)
	if err != nil {
		return nil, err
	}
	buf := make([]byte, reader.Len())
	if _, err = reader.Read(buf); err != nil {
		return nil, err
	}
	return parseTupleFromPageBuffer(buf)
}
