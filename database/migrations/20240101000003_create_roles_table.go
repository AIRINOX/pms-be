package migrations

import (
	"github.com/goravel/framework/contracts/database/schema"
	"github.com/goravel/framework/facades"
)

type M20240101000003CreateRolesTable struct{}

// Signature The unique signature for the migration.
func (r *M20240101000003CreateRolesTable) Signature() string {
	return "20240101000003_create_roles_table"
}

// Up Run the migrations.
func (r *M20240101000003CreateRolesTable) Up() error {
	return facades.Schema().Create("roles", func(table schema.Blueprint) {
		table.ID("id")
		table.String("key", 50)
		table.Unique("key")
		table.String("title", 100)
		table.Integer("order_index")
		table.TimestampsTz()

		table.Index("key")
		table.Index("order_index")
	})
}

// Down Reverse the migrations.
func (r *M20240101000003CreateRolesTable) Down() error {
	return facades.Schema().DropIfExists("roles")
}
