package migrations

import (
	"github.com/goravel/framework/contracts/database/schema"
	"github.com/goravel/framework/facades"
)

type M20240101000009CreateProductsTable struct{}

// Signature The unique signature for the migration.
func (r *M20240101000009CreateProductsTable) Signature() string {
	return "20240101000009_create_products_table"
}

// Up Run the migrations.
func (r *M20240101000009CreateProductsTable) Up() error {
	return facades.Schema().Create("products", func(table schema.Blueprint) {
		table.ID("id")
		table.String("title", 255)
		table.Text("description").Nullable()
		table.String("sku", 100).Nullable()
		table.Unique("sku")
		table.Boolean("is_raw_material").Default(false)
		table.UnsignedBigInteger("category_id").Nullable()
		table.UnsignedBigInteger("location_id").Nullable()
		table.Decimal("prix_achat").Nullable()
		table.Decimal("prix_vente").Nullable()
		table.String("unit", 50).Nullable()
		table.String("image_url", 255).Nullable()
		table.TimestampsTz()

		table.Foreign("category_id").References("id").On("categories")
		table.Foreign("location_id").References("id").On("storage_locations")
		table.Index("title")
		table.Index("sku")
		table.Index("is_raw_material")
		table.Index("category_id")
		table.Index("location_id")
	})
}

// Down Reverse the migrations.
func (r *M20240101000009CreateProductsTable) Down() error {
	return facades.Schema().DropIfExists("products")
}
