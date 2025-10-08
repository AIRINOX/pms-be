package migrations

import (
	"github.com/goravel/framework/contracts/database/schema"
	"github.com/goravel/framework/facades"
)

type M20240101000018CreateRecipeProductItemsTable struct{}

// Signature The unique signature for the migration.
func (r *M20240101000018CreateRecipeProductItemsTable) Signature() string {
	return "20240101000018_create_recipe_product_items_table"
}

// Up Run the migrations.
func (r *M20240101000018CreateRecipeProductItemsTable) Up() error {
	return facades.Schema().Create("recipe_product_items", func(table schema.Blueprint) {
		table.ID("id")
		table.UnsignedBigInteger("recipe_product_id")
		table.UnsignedBigInteger("material_product_id")
		table.Text("notes").Nullable()
		table.TimestampsTz()

		table.Foreign("recipe_product_id").References("id").On("recipe_products")
		table.Foreign("material_product_id").References("id").On("products")

		table.Index("recipe_product_id")
		table.Index("material_product_id")
	})
}

// Down Reverse the migrations.
func (r *M20240101000018CreateRecipeProductItemsTable) Down() error {
	return facades.Schema().DropIfExists("recipe_product_items")
}
