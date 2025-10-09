package migrations

import (
	"github.com/goravel/framework/contracts/database/schema"
	"github.com/goravel/framework/facades"
)

type M20240101000013CreateProductImagesTable struct{}

// Signature The unique signature for the migration.
func (r *M20240101000013CreateProductImagesTable) Signature() string {
	return "20240101000013_create_product_images_table"
}

// Up Run the migrations.
func (r *M20240101000013CreateProductImagesTable) Up() error {
	return facades.Schema().Create("product_images", func(table schema.Blueprint) {
		table.ID("id")
		table.UnsignedBigInteger("product_id")
		table.String("file_url", 255)
		table.String("file_name", 255)
		table.Integer("image_index")
		table.Boolean("is_primary").Default(false)
		table.TimestampsTz()

		table.Foreign("product_id").References("id").On("products")
		table.Index("product_id")
		table.Index("image_index")
		table.Index("is_primary")
	})
}

// Down Reverse the migrations.
func (r *M20240101000013CreateProductImagesTable) Down() error {
	return facades.Schema().DropIfExists("product_images")
}
