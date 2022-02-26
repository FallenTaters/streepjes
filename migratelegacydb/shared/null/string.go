package null

import (
	"database/sql"
	"encoding/json"
)

type String sql.NullString

func (i *String) UnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, &i.String)
}

func NewString(v string) String {
	return String(sql.NullString{
		Valid:  true,
		String: v,
	})
}
