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
		table.Foreign("product_id").References("id").On("products")
		table.Index("product_id")
		table.Index("title")

		table.ID("id")
		table.UnsignedBigInteger("product_id")
		table.UnsignedBigInteger("material_product_id")
		table.Text("notes").Nullable()
		table.TimestampsTz()

		table.Foreign("product_id").References("id").On("products")
		table.Foreign("material_product_id").References("id").On("products")

		table.Index("product_id")
		table.Index("material_product_id")
	})
}

// Down Reverse the migrations.
func (r *M20240101000014CreateRecipeProductsTable) Down() error {
	return facades.Schema().DropIfExists("recipe_products")
}
