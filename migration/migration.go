package migration

import (
	"crypto/md5"
	"fmt"
)

type migration struct {
	ID  string
	SQL string
}

func (m *migration) MD5() string {
	return fmt.Sprintf(
		"%x",
		md5.Sum([]byte(m.SQL)),
	)
}
