package null

import (
	"database/sql"
	"encoding/json"
	"strconv"
)

type Int sql.NullInt64

func (i Int) Int() int {
	return int(i.Int64)
}

func (i Int) String() string {
	return strconv.Itoa(i.Int())
}

func (i *Int) UnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, &i.Int64)
}

func NewInt(i int) Int {
	return Int(sql.NullInt64{
		Valid: true,
		Int64: int64(i),
	})
}
