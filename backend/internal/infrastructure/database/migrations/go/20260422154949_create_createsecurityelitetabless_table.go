package migrations

import (
	"github.com/kodia-studio/kodia/pkg/database"
)

// Migration_20260422154949 handles the creation of the createsecurityelitetabless table.
type Migration_20260422154949 struct{}

func (m *Migration_20260422154949) Up(schema *database.Schema) error {
	return schema.Create("createsecurityelitetabless", func(table *database.Blueprint) {
		table.ID()
		table.String("name")
		table.Timestamps()
		table.SoftDeletes()
	})
}

func (m *Migration_20260422154949) Down(schema *database.Schema) error {
	return schema.Drop("createsecurityelitetabless")
}
