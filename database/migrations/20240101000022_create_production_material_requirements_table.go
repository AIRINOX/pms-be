package migrations

import (
	"github.com/goravel/framework/contracts/database/schema"
	"github.com/goravel/framework/facades"
)

type M20240101000022CreateProductionMaterialRequirementsTable struct{}

// Signature The unique signature for the migration.
func (r *M20240101000022CreateProductionMaterialRequirementsTable) Signature() string {
	return "20240101000022_create_production_material_requirements_table"
}

// Up Run the migrations.
func (r *M20240101000022CreateProductionMaterialRequirementsTable) Up() error {
	return facades.Schema().Create("production_material_requirements", func(table schema.Blueprint) {
		table.ID("id")
		table.UnsignedBigInteger("order_fabrication_id")
		table.UnsignedBigInteger("material_variant_id")
		table.Decimal("required_quantity")
		table.Decimal("stock_quantity").Default(0)
		table.Decimal("request_quantity").Default(0)
		table.String("unit", 50).Nullable()
		table.Enum("status", []any{"pending", "requested", "available", "consumed"}).Nullable()
		table.TimestampsTz()

		table.Foreign("order_fabrication_id").References("id").On("order_fabrications")
		table.Foreign("material_variant_id").References("id").On("article_variants")

		table.Index("order_fabrication_id")
		table.Index("material_variant_id")
		table.Index("status")
	})
}

// Down Reverse the migrations.
func (r *M20240101000022CreateProductionMaterialRequirementsTable) Down() error {
	return facades.Schema().DropIfExists("production_material_requirements")
}
