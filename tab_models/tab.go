package tab_models

import "github.com/joshua0x/table_data_compare/db"

var tablePool = map[string]db.Tabler{}

func init() {
}

func GetTable(tname string) db.Tabler {
	return tablePool[tname]
}

