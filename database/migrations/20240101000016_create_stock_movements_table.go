package migrations

import (
	"github.com/goravel/framework/contracts/database/schema"
	"github.com/goravel/framework/facades"
)

type M20240101000016CreateStockMovementsTable struct{}

// Signature The unique signature for the migration.
func (r *M20240101000016CreateStockMovementsTable) Signature() string {
	return "20240101000016_create_stock_movements_table"
}

// Up Run the migrations.
func (r *M20240101000016CreateStockMovementsTable) Up() error {
	return facades.Schema().Create("stock_movements", func(table schema.Blueprint) {
		table.ID("id")
		table.UnsignedBigInteger("product_id").Nullable()
		table.UnsignedBigInteger("variant_id").Nullable()
		table.UnsignedBigInteger("location_id")
		table.Enum("movement_type", []any{"in", "out", "adjustment"})
		table.Decimal("quantity")
		table.String("unit", 50).Nullable()
		table.String("reference_type", 50).Nullable()
		table.UnsignedBigInteger("reference_id").Nullable()
		table.Text("notes").Nullable()
		table.UnsignedBigInteger("created_by")
		table.Timestamp("created_at").UseCurrent()

		table.Foreign("product_id").References("id").On("products")
		table.Foreign("variant_id").References("id").On("product_variants")
		table.Foreign("location_id").References("id").On("storage_locations")
		table.Foreign("created_by").References("id").On("users")
		table.Index("product_id")
		table.Index("variant_id")
		table.Index("location_id")
		table.Index("movement_type")
		table.Index("created_at")
	})
}

// Down Reverse the migrations.
func (r *M20240101000016CreateStockMovementsTable) Down() error {
	return facades.Schema().DropIfExists("stock_movements")
}
