package constants

import (
	"os"
	"path"
)

var cwd, _ = os.Getwd()

// update the path below
var DATA_FILE_PATH string = path.Join(cwd, "_data")

const (
	PAGE_SIZE            uint64 = 8 * 1024 // 8kb
	PAGE_VERSION                = 1.0
	FILE_HEADER_SIZE     uint64 = 2  // 2 bytes
	ITEM_ID_HEADER_SIZE  uint64 = 4  // 4 bytes
	PAGE_HEADER_SIZE     uint64 = 17 // 17 bytes
	ITEM_HEADER_SIZE     uint64 = 4  // 4 bytes
	TUPLE_HEADER_SIZE    uint64 = 6  // 6 bytes
	SPECIAL_SPACE_SIZE   uint64 = 8  // 8 bytes
	SPECIAL_SPACE_HEAD   uint32 = uint32(PAGE_SIZE - SPECIAL_SPACE_SIZE)
	FREE_SPACE_TAIL      uint32 = SPECIAL_SPACE_HEAD - 1
	FREE_SPACE_HEAD      uint64 = PAGE_HEADER_SIZE
	FREE_SPACE_PAGE_INIT uint64 = PAGE_SIZE - FILE_HEADER_SIZE - PAGE_HEADER_SIZE - SPECIAL_SPACE_SIZE
	STRING_HEADER_SIZE   uint64 = 4  // 4 byte
	INT_SIZE             int  = 8  // 8 byte
	BOOL_SIZE            int    = 1  // 1 byte
	TIMESTAMP_SIZE       int    = 15 // 15 byte
	FILE_MODE            int    = 0644
)
