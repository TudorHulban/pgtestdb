package migration

import "sort"

type migrations []migration

func (migrations migrations) SortByID() {
	sort.Slice(
		migrations,
		func(i, j int) bool {
			return migrations[i].ID < migrations[j].ID
		},
	)
}
