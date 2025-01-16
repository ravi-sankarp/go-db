package internal

import (
	"bytes"
	"fmt"
	"go-db/core/utils"
	"testing"
)

const dbName = "test"
const tableName = "users"

func TestSerializeAndDeserialize(t *testing.T) {
	var input = fileHeader{Last_inserted_page_no: 20}
	inputBuffer := new(bytes.Buffer)
	utils.Serialize(input, inputBuffer)
	var output fileHeader
	utils.DeSerialize(inputBuffer.Bytes(), &output, nil)
	fmt.Println(output.Last_inserted_page_no)
	if input.Last_inserted_page_no != output.Last_inserted_page_no {
		t.Error("Values don't match, input:", input, " output: ", output)
	}
	fmt.Println("SUCCESS")
}

func TestTableInitByteArrayDeserialization(t *testing.T) {
	byteArr := getPageInitData(nil)
	pageReader := bytes.NewReader(byteArr[2:])
	if _, err := parsePageHeadersFromBuffer(pageReader); err != nil {
		t.Fatal(err)
	} else {
		fmt.Println("SUCCESS")
	}
}
func TestTableCreation(t *testing.T) {
	if err := CreateDb(dbName); err != nil {
		t.Error(err)
		return
	}
	fmt.Println("Db created")
	if err := CreateTable(dbName, tableName); err != nil {
		t.Error(err)
		return
	}
	fmt.Println("Table created")
	fmt.Println("SUCCESS")
}

func TestInsertToTable(t *testing.T) {

	if err := InsertToTable(dbName, tableName, SAMPLE_TUPLE_INSERT); err != nil {
		t.Error(err)
		return

	}

	fmt.Println("SUCCESS")
}

func TestReadFromTable(t *testing.T) {

	result, err := ReadFromTable(dbName, tableName)
	if err != nil {
		t.Error(err)
		return
	}

	actualResultSet := make([]SAMPLE_TABLE_STRUCT, len(result))
	for i, _ := range result {
		var row SAMPLE_TABLE_STRUCT
		if err := result[i].Scan(&row.Name, &row.Email, &row.Age, &row.Deleted, &row.Created_on); err != nil {
			t.Error(err)
			return
		}
		actualResultSet[i] = row
	}
	fmt.Println(actualResultSet)
	fmt.Println("SUCCESS")
}
