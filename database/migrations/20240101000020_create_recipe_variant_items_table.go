package migrations

import (
	"github.com/goravel/framework/contracts/database/schema"
	"github.com/goravel/framework/facades"
)

type M20240101000020CreateRecipeVariantItemsTable struct{}

// Signature The unique signature for the migration.
func (r *M20240101000020CreateRecipeVariantItemsTable) Signature() string {
	return "20240101000020_create_recipe_variant_items_table"
}

// Up Run the migrations.
func (r *M20240101000020CreateRecipeVariantItemsTable) Up() error {
	return facades.Schema().Create("recipe_variant_items", func(table schema.Blueprint) {
		table.ID("id")
		table.UnsignedBigInteger("recipe_variant_id")
		table.UnsignedBigInteger("material_variant_id")
		table.Decimal("quantity")
		table.Text("notes").Nullable()
		table.TimestampsTz()

		table.Foreign("recipe_variant_id").References("id").On("recipe_variants")
		table.Foreign("material_variant_id").References("id").On("article_variants")

		table.Index("recipe_variant_id")
		table.Index("material_variant_id")
	})
}

// Down Reverse the migrations.
func (r *M20240101000020CreateRecipeVariantItemsTable) Down() error {
	return facades.Schema().DropIfExists("recipe_variant_items")
}
