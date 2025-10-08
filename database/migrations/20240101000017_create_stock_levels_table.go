package migrations

import (
	"github.com/goravel/framework/contracts/database/schema"
	"github.com/goravel/framework/facades"
)

type M20240101000017CreateStockLevelsTable struct{}

// Signature The unique signature for the migration.
func (r *M20240101000017CreateStockLevelsTable) Signature() string {
	return "20240101000017_create_stock_levels_table"
}

// Up Run the migrations.
func (r *M20240101000017CreateStockLevelsTable) Up() error {
	return facades.Schema().Create("stock_levels", func(table schema.Blueprint) {
		table.ID("id")
		table.UnsignedBigInteger("product_id").Nullable()
		table.UnsignedBigInteger("variant_id").Nullable()
		table.UnsignedBigInteger("location_id")
		table.Decimal("quantity").Default(0)
		table.String("unit", 50).Nullable()
		table.Timestamp("last_updated").Nullable()

		table.Foreign("product_id").References("id").On("products")
		table.Foreign("variant_id").References("id").On("product_variants")
		table.Foreign("location_id").References("id").On("storage_locations")
		table.Unique("product_id", "variant_id", "location_id")
		table.Index("location_id")
		table.Index("quantity")
	})
}

// Down Reverse the migrations.
func (r *M20240101000017CreateStockLevelsTable) Down() error {
	return facades.Schema().DropIfExists("stock_levels")
}
