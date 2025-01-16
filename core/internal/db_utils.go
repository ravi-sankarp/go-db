package internal

import (
	"go-db/core/constants"
	"path"
)

func getTableFilePath(dbName string, tableName string) string {
	return path.Join(constants.DATA_FILE_PATH, dbName, tableName)
}
