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
		table.UnsignedBigInteger("article_id")
		table.UnsignedBigInteger("variant_id")
		table.Decimal("output_quantity")
		table.Text("notes").Nullable()
		table.TimestampsTz()

		table.Foreign("article_id").References("id").On("articles")
		table.Foreign("variant_id").References("id").On("article_variants")

		table.Index("article_id")
		table.Index("variant_id")
		table.Unique("article_id", "variant_id")
	})
}

// Down Reverse the migrations.
func (r *M20240101000019CreateRecipeVariantsTable) Down() error {
	return facades.Schema().DropIfExists("recipe_variants")
}
