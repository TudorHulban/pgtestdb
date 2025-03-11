package pgtestdb

import (
	"fmt"
	"strings"
)

func updateDBNameInConnection(connectionString, newDBName string) string {
	parts := strings.Split(connectionString, "/")
	if len(parts) > 3 {
		parts[len(parts)-1] = fmt.Sprintf(
			"%s?",
			newDBName,
		)
	}

	return strings.Join(parts, "/")
}
