package internal

import (
	"go-db/core/constants"
)

func calcPageOffset(pageNo uint16) int64 {
	return int64(pageNo)*int64(constants.PAGE_SIZE) + int64(constants.FILE_HEADER_SIZE)
}

func calcItemHeaderOffset(tupleCount uint8, pageNo uint16) int64 {
	return calcPageOffset(pageNo) + int64(constants.PAGE_HEADER_SIZE) + (int64(tupleCount) * int64(constants.ITEM_HEADER_SIZE))
}

func calcTupleOffset(pageNo uint16, offsetInPage uint16) int64 {
	return calcPageOffset(pageNo) + int64(offsetInPage)
}

func checkIfNewPageIsRequired(tupleCount uint8, bufLen int, freeSpaceHead uint32, freeSpaceTail uint32) bool {
	return int(tupleCount)*int(constants.ITEM_HEADER_SIZE)+bufLen+int(freeSpaceHead) >= int(freeSpaceTail)
}
