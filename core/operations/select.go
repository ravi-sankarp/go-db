package operations

import (
	"fmt"
	"os"
	"strconv"
)

type test struct {
	id   uint16
	name string
}

func Read(query string) any {
	return "result"
}

func Write() any {
	val := test{id: 12, name: "123"}
	err := os.WriteFile("_data/table", []byte("id,name:"+strconv.Itoa(int(val.id))+val.name), os.ModeDevice)
	fmt.Print(err)
	return nil
}
