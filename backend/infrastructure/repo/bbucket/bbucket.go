package bbucket

import (
	"encoding/binary"
)

var (
	memberBucket = []byte("members")
	orderBucket  = []byte("orders")
)

func itob(v int) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(v))

	return b
}
