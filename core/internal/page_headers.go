package internal

import "time"

type fileHeader struct {
	Last_inserted_page_no uint16
}

type pageHeader struct {
	Free_space_head    uint32 // for referncing the head of the free space
	Free_space_tail    uint32
	Tuple_count        uint8 // for keeping a count of no of tuples in the page
	Special_space_tail uint32
	Version            float32
}

type itemHeader struct {
	Byte_offset uint16
	Length      uint16
}

type itemHeaders []itemHeader

type stringHeader struct {
	length uint32
}
type tupleHeader struct {
	Ct_id      uint32
	Attributes uint8
}

type Timestamp time.Time

type SampleTable struct {
	Id         int
	Name       string
	Email      string
	Created_on Timestamp
	Deleted    bool
}
