package migrations

import (
	"github.com/goravel/framework/contracts/database/schema"
	"github.com/goravel/framework/facades"
)

type M20240101000012CreateProductVariantsTable struct{}

// Signature The unique signature for the migration.
func (r *M20240101000012CreateProductVariantsTable) Signature() string {
	return "20240101000012_create_product_variants_table"
}

// Up Run the migrations.
func (r *M20240101000012CreateProductVariantsTable) Up() error {
	return facades.Schema().Create("product_variants", func(table schema.Blueprint) {
		table.ID("id")
		table.UnsignedBigInteger("product_id")
		table.String("title", 255).Nullable()
		table.Text("description").Nullable()
		table.String("sku", 100).Nullable()
		table.Json("attributes")
		table.Decimal("prix_achat").Nullable()
		table.Decimal("prix_vente").Nullable()
		table.String("unit", 50).Nullable()
		table.String("image_url", 255).Nullable()
		table.Integer("image_index").Nullable()
		table.Boolean("is_active").Default(true)
		table.TimestampsTz()

		table.Foreign("product_id").References("id").On("products")
		table.Index("product_id")
		table.Index("sku")
		table.Index("is_active")
	})
}

// Down Reverse the migrations.
func (r *M20240101000012CreateProductVariantsTable) Down() error {
	return facades.Schema().DropIfExists("product_variants")
}
