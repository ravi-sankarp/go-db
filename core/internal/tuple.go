package internal

import (
	dbErrors "go-db/core/errors"
	"reflect"
	"time"
)

type attribute struct {
	name     string
	value    reflect.Value
	dataType reflect.Kind
}

type Tuple []attribute

func convertToTargetType(value reflect.Value, target any) error {
	switch target.(type) {
	case *int:
	case *int64:
		*target.(*int64) = value.Interface().(int64)
	case *string:
		*target.(*string) = value.Interface().(string)
	case *bool:
		*target.(*bool) = value.Interface().(bool)
	case *time.Time:
		*target.(*time.Time) = value.Interface().(time.Time)
	default:
		return dbErrors.NewDbError("Invalid type for target variable")
	}
	return nil
}

func (tuple *Tuple) Scan(dest ...any) error {
	for i := range dest {
		convertToTargetType((*tuple)[i].value, dest[i])
	}
	return nil
}
