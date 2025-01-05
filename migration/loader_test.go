package migration

import (
	"regexp"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRegex(t *testing.T) {
	regex := `^V\d{4}_[a-zA-Z0-9]+\.sql$`

	regexObject := regexp.MustCompile(regex)

	require.True(t,
		regexObject.MatchString(
			"V0003_differentnaming.sql",
		),
	)
}
