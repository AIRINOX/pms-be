package migrations

import (
	"github.com/goravel/framework/contracts/database/schema"
	"github.com/goravel/framework/facades"
)

type M20240101000010CreateProductAttributesTable struct{}

// Signature The unique signature for the migration.
func (r *M20240101000010CreateProductAttributesTable) Signature() string {
	return "20240101000010_create_product_attributes_table"
}

// Up Run the migrations.
func (r *M20240101000010CreateProductAttributesTable) Up() error {
	return facades.Schema().Create("product_attributes", func(table schema.Blueprint) {
		table.ID("id")
		table.UnsignedBigInteger("product_id")
		table.String("key", 100)
		table.String("title", 255)
		table.Integer("order_index").Default(0)
		table.Timestamp("created_at").UseCurrent()

		table.Foreign("product_id").References("id").On("products")
		table.Unique("product_id", "key")
		table.Index("product_id")
		table.Index("order_index")
	})
}

// Down Reverse the migrations.
func (r *M20240101000010CreateProductAttributesTable) Down() error {
	return facades.Schema().DropIfExists("product_attributes")
}
