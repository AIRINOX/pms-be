package migrations

import (
	"github.com/goravel/framework/contracts/database/schema"
	"github.com/goravel/framework/facades"
)

type M20240101000023CreateStockRequestsTable struct{}

// Signature The unique signature for the migration.
func (r *M20240101000023CreateStockRequestsTable) Signature() string {
	return "20240101000023_create_stock_requests_table"
}

// Up Run the migrations.
func (r *M20240101000023CreateStockRequestsTable) Up() error {
	return facades.Schema().Create("stock_requests", func(table schema.Blueprint) {
		table.ID("id")
		table.UnsignedBigInteger("order_fabrication_id")
		table.UnsignedBigInteger("material_variant_id")
		table.Decimal("requested_quantity")
		table.String("unit", 50).Nullable()
		table.Enum("status", []any{"pending", "approved", "rejected", "fulfilled"}).Nullable()
		table.UnsignedBigInteger("requested_by")
		table.Timestamp("requested_at").Nullable()
		table.Timestamp("fulfilled_at").Nullable()
		table.Text("notes").Nullable()

		table.Foreign("order_fabrication_id").References("id").On("order_fabrications")
		table.Foreign("material_variant_id").References("id").On("product_variants")
		table.Foreign("requested_by").References("id").On("users")

		table.Index("order_fabrication_id")
		table.Index("material_variant_id")
		table.Index("requested_by")
		table.Index("status")
		table.Index("requested_at")
	})
}

// Down Reverse the migrations.
func (r *M20240101000023CreateStockRequestsTable) Down() error {
	return facades.Schema().DropIfExists("stock_requests")
}
