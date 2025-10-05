package migrations

import (
	"github.com/goravel/framework/contracts/database/schema"
	"github.com/goravel/framework/facades"
)

type M20240101000002CreateCategoriesTable struct{}

// Signature The unique signature for the migration.
func (r *M20240101000002CreateCategoriesTable) Signature() string {
	return "20240101000002_create_categories_table"
}

// Up Run the migrations.
func (r *M20240101000002CreateCategoriesTable) Up() error {
	return facades.Schema().Create("categories", func(table schema.Blueprint) {
		table.ID("id")
		table.String("title", 255)
		table.Text("description").Nullable()
		table.UnsignedBigInteger("parent_id").Nullable()
		table.TimestampsTz()

		table.Foreign("parent_id").References("id").On("categories")
		table.Index("parent_id")
	})
}

// Down Reverse the migrations.
func (r *M20240101000002CreateCategoriesTable) Down() error {
	return facades.Schema().DropIfExists("categories")
}
