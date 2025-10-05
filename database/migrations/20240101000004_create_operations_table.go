package migrations

import (
	"github.com/goravel/framework/contracts/database/schema"
	"github.com/goravel/framework/facades"
)

type M20240101000004CreateOperationsTable struct{}

// Signature The unique signature for the migration.
func (r *M20240101000004CreateOperationsTable) Signature() string {
	return "20240101000004_create_operations_table"
}

// Up Run the migrations.
func (r *M20240101000004CreateOperationsTable) Up() error {
	return facades.Schema().Create("operations", func(table schema.Blueprint) {
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
func (r *M20240101000004CreateOperationsTable) Down() error {
	return facades.Schema().DropIfExists("operations")
}
