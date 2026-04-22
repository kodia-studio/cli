package migrations

import (
	"github.com/kodia-studio/kodia/pkg/database"
)

// Migration_20260422150028 handles the creation of the createwebhookhistoriestables table.
type Migration_20260422150028 struct{}

func (m *Migration_20260422150028) Up(schema *database.Schema) error {
	return schema.Create("createwebhookhistoriestables", func(table *database.Blueprint) {
		table.ID()
		table.String("name")
		table.Timestamps()
		table.SoftDeletes()
	})
}

func (m *Migration_20260422150028) Down(schema *database.Schema) error {
	return schema.Drop("createwebhookhistoriestables")
}
