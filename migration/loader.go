package migration

import (
	"fmt"
	"io/fs"
	"regexp"
)

func load(directory fs.FS, regex *regexp.Regexp) (migrations, error) {
	var result []migration

	if err := fs.WalkDir(
		directory,
		".",
		func(path string, d fs.DirEntry, errReadFile error) error {
			if errReadFile != nil {
				return errReadFile
			}

			if !regex.MatchString(path) {
				return nil
			}

			content, errReadContent := fs.ReadFile(directory, path)
			if errReadContent != nil {
				return errReadContent
			}

			result = append(
				result,
				migration{
					ID:  d.Name(),
					SQL: string(content),
				},
			)

			return nil
		},
	); err != nil {
		return nil,
			fmt.Errorf("load: %w", err)
	}

	return result, nil
}
