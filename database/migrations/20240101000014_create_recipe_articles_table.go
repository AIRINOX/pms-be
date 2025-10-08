package migrations

import (
	"github.com/goravel/framework/contracts/database/schema"
	"github.com/goravel/framework/facades"
)

type M20240101000014CreateRecipeProductsTable struct{}

// Signature The unique signature for the migration.
func (r *M20240101000014CreateRecipeProductsTable) Signature() string {
	return "20240101000014_create_recipe_products_table"
}

// Up Run the migrations.
func (r *M20240101000014CreateRecipeProductsTable) Up() error {
	return facades.Schema().Create("recipe_products", func(table schema.Blueprint) {
		table.ID("id")
		table.String("title", 255)
		table.Text("description").Nullable()
		table.UnsignedBigInteger("product_id")
		table.TimestampsTz()

		table.Foreign("product_id").References("id").On("products")
		table.Index("product_id")
		table.Index("title")
	})
}

// Down Reverse the migrations.
func (r *M20240101000014CreateRecipeProductsTable) Down() error {
	return facades.Schema().DropIfExists("recipe_products")
}
