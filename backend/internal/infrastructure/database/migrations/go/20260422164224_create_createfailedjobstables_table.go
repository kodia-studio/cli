package migrations

import (
	"github.com/kodia-studio/kodia/pkg/database"
)

// Migration_20260422164224 handles the creation of the createfailedjobstables table.
type Migration_20260422164224 struct{}

func (m *Migration_20260422164224) Up(schema *database.Schema) error {
	return schema.Create("createfailedjobstables", func(table *database.Blueprint) {
		table.ID()
		table.String("name")
		table.Timestamps()
		table.SoftDeletes()
	})
}

func (m *Migration_20260422164224) Down(schema *database.Schema) error {
	return schema.Drop("createfailedjobstables")
}
