package migrations

import (
	"github.com/goravel/framework/contracts/database/schema"
	"github.com/goravel/framework/facades"
)

type M20240101000006CreateStorageLocationsTable struct{}

// Signature The unique signature for the migration.
func (r *M20240101000006CreateStorageLocationsTable) Signature() string {
	return "20240101000006_create_storage_locations_table"
}

// Up Run the migrations.
func (r *M20240101000006CreateStorageLocationsTable) Up() error {
	return facades.Schema().Create("storage_locations", func(table schema.Blueprint) {
		table.ID("id")
		table.String("name", 255)
		table.Text("description").Nullable()
		table.TimestampsTz()

		table.Index("name")
	})
}

// Down Reverse the migrations.
func (r *M20240101000006CreateStorageLocationsTable) Down() error {
	return facades.Schema().DropIfExists("storage_locations")
}
