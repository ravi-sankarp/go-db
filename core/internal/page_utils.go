package internal

import (
	"bytes"
	"fmt"
	"go-db/core/constants"
	dbErrors "go-db/core/errors"
	"go-db/core/utils"
	"io"
	"os"
)

const (
	lastInsertedPageNoOffset = 0
)

func ensurePageSize(buffer *bytes.Buffer, targetSize int) error {
	currentLen := buffer.Len()

	if currentLen > targetSize {
		panic("Buffer Size overflow")
	} else if currentLen < targetSize {
		// Pad the buffer with zeros
		padding := make([]byte, targetSize-currentLen)
		buffer.Write(padding)
	}

	return nil
}

func appendAtFixedOffset(buffer *bytes.Buffer, byteOffset int, data []byte) error {
	// Overwrite the contents at the specified offset
	rawBuffer := buffer.Bytes()
	copy(rawBuffer[byteOffset:], data)

	return nil
}

func getPageInitData(initData *bytes.Buffer) []byte {
	targetSize := int(constants.PAGE_SIZE)
	pageHeader := pageHeader{
		Special_space_tail: uint32(constants.PAGE_SIZE - constants.SPECIAL_SPACE_SIZE),
		Version:            constants.PAGE_VERSION,
		Free_space_head:    uint32(constants.FREE_SPACE_HEAD),
		Free_space_tail:    constants.FREE_SPACE_TAIL,
		Tuple_count:        0,
	}
	var itemHeader itemHeader
	var bufLen int
	if initData != nil {
		bufLen = initData.Len()
		pageHeader.Free_space_tail -= uint32(bufLen - 1)
		itemHeader.Byte_offset = uint16(pageHeader.Free_space_tail + 1)
		itemHeader.Length = uint16(bufLen)
		pageHeader.Tuple_count += 1
	}
	inputBuffer := new(bytes.Buffer)
	if initData == nil { // fileheader is only needed once per file
		utils.Serialize(fileHeader{Last_inserted_page_no: 0}, inputBuffer)
		targetSize += int(constants.FILE_HEADER_SIZE)
	}
	utils.Serialize(pageHeader, inputBuffer)
	// ensurePageSize(inputBuffer, targetSize)
	if initData != nil {
		itemHeaderBuf := new(bytes.Buffer)
		utils.Serialize(itemHeader, itemHeaderBuf)
		initBytes := initData.Bytes()
		itemOffset := int(constants.PAGE_HEADER_SIZE) + int(constants.FILE_HEADER_SIZE)
		tupleOffset := int(pageHeader.Free_space_tail+1) + int(constants.FILE_HEADER_SIZE)
		appendAtFixedOffset(inputBuffer, itemOffset, itemHeaderBuf.Bytes()) // item header
		appendAtFixedOffset(inputBuffer, tupleOffset, initBytes)            // tuple
	}
	return inputBuffer.Bytes()
}

func updateLastInsertedPageNo(file *os.File, pageNo uint16) error {
	buf := new(bytes.Buffer)
	utils.Serialize(fileHeader{Last_inserted_page_no: pageNo}, buf)
	_, err := file.WriteAt(buf.Bytes(), lastInsertedPageNoOffset)
	return err

}

func flushUpdatedPageAndItemHeader(file *os.File, pageHeader pageHeader, itemHeader itemHeader, pageNo uint16) error {
	var (
		pageBuf       = new(bytes.Buffer)
		pageBufOffset = calcPageOffset(pageNo)
		itemBuf       = new(bytes.Buffer)
		itemBufOffset = calcItemHeaderOffset(pageHeader.Tuple_count, pageNo)
	)
	utils.Serialize(pageHeader, pageBuf)
	utils.Serialize(itemHeader, itemBuf)

	if _, err := file.WriteAt(pageBuf.Bytes(), pageBufOffset); err != nil {
		return err
	}
	_, err := file.WriteAt(itemBuf.Bytes(), itemBufOffset)
	return err
}

func flushTupleToDisk(file *os.File, data []byte, offset int64) error {
	_, err := file.WriteAt(data, offset)
	return err
}

func appendPage(page []byte, file *os.File) error {
	if _, err := file.Write(page); err != nil {
		return err
	}
	return nil
}

func readFromFileOffset(file *os.File, bufferSize int, offset int64) ([]byte, error) {
	buffer := make([]byte, bufferSize)
	if _, err := file.Seek(offset, io.SeekStart); err != nil {
		fmt.Printf("Error reading file: %v\n", err)
		return buffer, dbErrors.NewDbError("Error reading data")
	}

	if _, err := file.Read(buffer); err != nil {
		fmt.Printf("Error reading file: %v\n", err)
		return buffer, dbErrors.NewDbError("Error reading data")
	}
	return buffer, nil
}

func readFromFileOffsetAndDeSerialize[T any](file *os.File, offset int64, bufferSize int, target T) (any, error) {
	buffer, err := readFromFileOffset(file, bufferSize, offset)
	if err != nil {
		return target, err
	}
	utils.DeSerialize(buffer, target, nil)
	return target, nil
}

func seekBufferReader(reader *bytes.Reader, offset int64, whence int) error {
	if _, err := reader.Seek(offset, whence); err != nil {
		fmt.Printf("Error reading buf: %v\n", err)
		return dbErrors.NewDbError("Error reading data")
	}
	return nil
}

func readFromBufferOffsetAndDeSerialize[T any](buf *bytes.Reader, offset int64, target T, whence int) error {
	err := seekBufferReader(buf, offset, whence)
	if err != nil {
		return err
	}
	utils.DeSerialize(nil, target, buf)
	return nil
}

func getFileHeaders(file *os.File) (fileHeader, error) {
	var (
		offset     int64 = (0)
		bufferSize       = int(constants.FILE_HEADER_SIZE)
		header     fileHeader
	)
	readFromFileOffsetAndDeSerialize(file, offset, bufferSize, &header)
	return header, nil
}

func getPageBuf(file *os.File, pageNo uint16) ([]byte, error) {
	var (
		offset     int64 = int64(int64(constants.FILE_HEADER_SIZE) + (int64(pageNo) * int64(constants.PAGE_SIZE)))
		bufferSize       = int(constants.PAGE_SIZE)
	)

	bytes, err := readFromFileOffset(file, bufferSize, offset)
	if err != nil {
		return bytes, err
	}
	return bytes, nil
}

func parsePageHeadersFromBuffer(pageBuf *bytes.Reader) (pageHeader, error) {
	var (
		header pageHeader
	)
	readFromBufferOffsetAndDeSerialize(pageBuf, 0, &header, 0)
	fmt.Println("head ", header.Free_space_head, " tail ", header.Free_space_tail)
	return header, nil
}

func parseItemHeadersFromBuffer(pageBuf *bytes.Reader, noOfItems uint8) (itemHeaders, error) {
	var (
		itemHeaders = make(itemHeaders, noOfItems)
	)
	var i int64 = 0
	for i = 0; i < int64(noOfItems); i++ {
		var (
			header itemHeader
			offset = int64(constants.PAGE_HEADER_SIZE) + i*int64(constants.ITEM_HEADER_SIZE)
		)
		if err := readFromBufferOffsetAndDeSerialize(pageBuf, offset, &header, 0); err != nil {
			return nil, err
		}
		fmt.Println(header)
		itemHeaders[i] = header
	}
	return itemHeaders, nil
}

func getTuplesFromPage(page []byte) ([]Tuple, error) {
	var (
		pageHeader  pageHeader
		itemHeaders itemHeaders
		err         error
	)
	pageReader := bytes.NewReader(page)
	if pageHeader, err = parsePageHeadersFromBuffer(pageReader); err != nil {
		return nil, err
	}
	if itemHeaders, err = parseItemHeadersFromBuffer(pageReader, pageHeader.Tuple_count); err != nil {
		return nil, err
	}
	fmt.Println(itemHeaders)
	result := make([]Tuple, pageHeader.Tuple_count)
	for i, item := range itemHeaders {
		if tuple, err := parseTuple(pageReader, item.Byte_offset); err != nil {
			return nil, err
		} else {
			result[i] = tuple
		}

	}
	return result, nil
}
