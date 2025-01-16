package operations

import (
	"go-db/core/internal"
)

func Insert(dbName string, tableName string, data internal.Tuple) error {
	err := internal.InsertToTable(dbName, tableName, data)
	return err
}
