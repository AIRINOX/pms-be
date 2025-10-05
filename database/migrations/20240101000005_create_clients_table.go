package migrations

import (
	"github.com/goravel/framework/contracts/database/schema"
	"github.com/goravel/framework/facades"
)

type M20240101000005CreateClientsTable struct{}

// Signature The unique signature for the migration.
func (r *M20240101000005CreateClientsTable) Signature() string {
	return "20240101000005_create_clients_table"
}

// Up Run the migrations.
func (r *M20240101000005CreateClientsTable) Up() error {
	return facades.Schema().Create("clients", func(table schema.Blueprint) {
		table.ID("id")
		table.String("name", 255)
		table.String("phone", 50).Nullable()
		table.String("email", 255).Nullable()
		table.Text("address").Nullable()
		table.TimestampsTz()

		table.Index("name")
		table.Index("email")
	})
}

// Down Reverse the migrations.
func (r *M20240101000005CreateClientsTable) Down() error {
	return facades.Schema().DropIfExists("clients")
}
