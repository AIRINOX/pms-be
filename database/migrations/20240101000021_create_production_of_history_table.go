package migrations

import (
	"github.com/goravel/framework/contracts/database/schema"
	"github.com/goravel/framework/facades"
)

type M20240101000021CreateProductionOfHistoryTable struct{}

// Signature The unique signature for the migration.
func (r *M20240101000021CreateProductionOfHistoryTable) Signature() string {
	return "20240101000021_create_production_of_history_table"
}

// Up Run the migrations.
func (r *M20240101000021CreateProductionOfHistoryTable) Up() error {
	return facades.Schema().Create("production_of_history", func(table schema.Blueprint) {
		table.ID("id")
		table.UnsignedBigInteger("order_fabrication_id")
		table.UnsignedBigInteger("operation_id")
		table.UnsignedBigInteger("user_id")
		table.String("status", 50)
		table.Text("notes").Nullable()
		table.Timestamp("status_at").UseCurrent()

		table.Foreign("order_fabrication_id").References("id").On("order_fabrications")
		table.Foreign("operation_id").References("id").On("operations")
		table.Foreign("user_id").References("id").On("users")

		table.Index("order_fabrication_id")
		table.Index("operation_id")
		table.Index("user_id")
		table.Index("status")
		table.Index("status_at")
	})
}

// Down Reverse the migrations.
func (r *M20240101000021CreateProductionOfHistoryTable) Down() error {
	return facades.Schema().DropIfExists("production_of_history")
}
