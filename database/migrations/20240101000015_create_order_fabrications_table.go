package migrations

import (
	"github.com/goravel/framework/contracts/database/schema"
	"github.com/goravel/framework/facades"
)

type M20240101000015CreateOrderFabricationsTable struct{}

// Signature The unique signature for the migration.
func (r *M20240101000015CreateOrderFabricationsTable) Signature() string {
	return "20240101000015_create_order_fabrications_table"
}

// Up Run the migrations.
func (r *M20240101000015CreateOrderFabricationsTable) Up() error {
	return facades.Schema().Create("order_fabrications", func(table schema.Blueprint) {
		table.ID("id")
		table.String("order_number", 100)
		table.Unique("order_number")
		table.UnsignedBigInteger("product_id")
		table.UnsignedBigInteger("variant_id").Nullable()
		table.Decimal("quantity")
		table.UnsignedBigInteger("client_id")
		table.UnsignedBigInteger("client_site_id").Nullable()
		table.Enum("status", []any{
			"pending_validation", "validated", "rejected", "material_requested",
			"ready_to_produce", "cutting_started", "cutting_paused", "cutting_completed",
			"folding_started", "folding_completed", "assembly_started", "assembly_completed",
			"finishing_started", "finishing_completed", "ready_for_delivery", "delivered", "cancelled",
		}).Default("pending_validation")
		table.Integer("priority").Default(0)
		table.Date("deadline_date").Nullable()
		table.Text("notes").Nullable()
		table.UnsignedBigInteger("created_by")
		table.TimestampsTz()

		table.Foreign("product_id").References("id").On("products")
		table.Foreign("variant_id").References("id").On("product_variants")
		table.Foreign("client_id").References("id").On("clients")
		table.Foreign("client_site_id").References("id").On("client_sites")
		table.Foreign("created_by").References("id").On("users")
		table.Index("order_number")
		table.Index("status")
		table.Index("priority")
		table.Index("deadline_date")
	})
}

// Down Reverse the migrations.
func (r *M20240101000015CreateOrderFabricationsTable) Down() error {
	return facades.Schema().DropIfExists("order_fabrications")
}
