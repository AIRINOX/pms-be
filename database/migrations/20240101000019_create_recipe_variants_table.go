package migrations

import (
	"github.com/goravel/framework/contracts/database/schema"
	"github.com/goravel/framework/facades"
)

type M20240101000019CreateRecipeVariantsTable struct{}

// Signature The unique signature for the migration.
func (r *M20240101000019CreateRecipeVariantsTable) Signature() string {
	return "20240101000019_create_recipe_variants_table"
}

// Up Run the migrations.
func (r *M20240101000019CreateRecipeVariantsTable) Up() error {
	return facades.Schema().Create("recipe_variants", func(table schema.Blueprint) {
		table.ID("id")
		table.UnsignedBigInteger("product_id")
		table.UnsignedBigInteger("variant_id")
		table.Decimal("output_quantity")
		table.Text("notes").Nullable()
		table.TimestampsTz()

		table.Foreign("product_id").References("id").On("products")
		table.Foreign("variant_id").References("id").On("product_variants")

		table.Index("product_id")
		table.Index("variant_id")
		table.Unique("product_id", "variant_id")
	})
}

// Down Reverse the migrations.
func (r *M20240101000019CreateRecipeVariantsTable) Down() error {
	return facades.Schema().DropIfExists("recipe_variants")
}
