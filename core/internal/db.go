package internal

import (
	"bytes"
	"errors"
	"go-db/core/constants"
	dbErrors "go-db/core/errors"
	"os"
	"path"
)

func createDataFolder() error {
	if _, err := os.Stat(constants.DATA_FILE_PATH); errors.Is(err, os.ErrNotExist) {
		error := os.Mkdir(constants.DATA_FILE_PATH, os.ModeDir)
		return error
	} else {
		return nil
	}
}

func getDbDirName(name string) string {
	return path.Join(constants.DATA_FILE_PATH, name)
}

func CreateDb(name string) error {
	err := createDataFolder()
	if err != nil {
		return err
	}
	dirName := getDbDirName(name)
	if _, err := os.Stat(dirName); errors.Is(err, os.ErrNotExist) {
		error := os.Mkdir(dirName, os.ModeDir)
		return error
	} else {
		return dbErrors.NewDbError(dbErrors.DUPLICATE_DB_NAME)
	}
}
func DeleteDb(name string) error {
	dirName := getDbDirName(name)
	if err := os.RemoveAll(dirName); errors.Is(err, os.ErrNotExist) {
		return err
	}
	return nil
}

func CreateTable(dbName string, tableName string) error {
	fileName := getTableFilePath(dbName, tableName)
	if _, err := os.Stat(fileName); errors.Is(err, os.ErrNotExist) {
		error := os.WriteFile(fileName, getPageInitData(nil), os.ModeAppend)
		return error
	} else {
		return dbErrors.NewDbError(dbErrors.DUPLICATE_DB_NAME)
	}
}

func validateTableNameAndDb(dbName string, tableName string) error {
	// TODO check if tablename and dbname is valid
	return nil
}

func openTableFile(dbName string, tableName string) (*os.File, error) {
	if err := validateTableNameAndDb(dbName, tableName); err != nil {
		return nil, err
	}
	file, err := os.OpenFile(getTableFilePath(dbName, tableName), os.O_RDWR, os.FileMode(constants.FILE_MODE))
	if err != nil {
		return nil, err
	}
	return file, nil
}

func InsertToTable(dbName string, tableName string, data Tuple) error {
	file, err := openTableFile(dbName, tableName)
	if err != nil {
		return err
	}
	defer file.Close()
	var (
		fileHeader fileHeader
		pageHeader pageHeader
		page       []byte
		dataBuf    *bytes.Buffer
	)

	if fileHeader, err = getFileHeaders(file); err != nil {
		return err
	}
	if page, err = getPageBuf(file, fileHeader.Last_inserted_page_no); err != nil {
		return err
	}
	pageReader := bytes.NewReader(page)
	if pageHeader, err = parsePageHeadersFromBuffer(pageReader); err != nil {
		return err
	}

	// TODO dynamic table schema

	if dataBuf, err = parseDataToBinaryTuple(data); err != nil {
		return err
	}

	bufLen := dataBuf.Len()
	if checkIfNewPageIsRequired(pageHeader.Tuple_count+1, bufLen, pageHeader.Free_space_head, pageHeader.Free_space_tail) {
		// panic(fmt.Sprint("New page, Tuple Count: ", pageHeader.Tuple_count))
		appendPage(getPageInitData(dataBuf), file) // append the new page
		calcAndUpdateLastInsertedPageNo(file, fileHeader)
		return nil
	}
	var itemHeader itemHeader
	updatePageAndItemHeaderFromBufferLength(&pageHeader, &itemHeader, bufLen)
	if err = flushUpdatedPageAndItemHeader(file, pageHeader, itemHeader, fileHeader.Last_inserted_page_no); err != nil {
		return err
	}
	tupleOffset := calcTupleOffset(fileHeader.Last_inserted_page_no, itemHeader.Byte_offset)
	flushTupleToDisk(file, dataBuf.Bytes(), tupleOffset)
	return nil
}

func ReadFromTable(dbName string, tableName string) ([]Tuple, error) {
	file, err := openTableFile(dbName, tableName)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var (
		fileHeader fileHeader
		page       []byte
	)

	if fileHeader, err = getFileHeaders(file); err != nil {
		return nil, err
	}
	var i uint16 = 0
	var result []Tuple
	for i = 0; i <= fileHeader.Last_inserted_page_no; i++ {
		if page, err = getPageBuf(file, i); err != nil {
			return nil, err
		}
		tuples, err := getTuplesFromPage(&page)
		if err != nil {
			return nil, err
		}
		result = append(result, tuples...)
	}

	return result, nil
}
