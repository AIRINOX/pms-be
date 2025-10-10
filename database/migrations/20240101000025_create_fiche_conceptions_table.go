package migrations

import (
	"github.com/goravel/framework/contracts/database/schema"
	"github.com/goravel/framework/facades"
)

type M20240101000025CreateFicheConceptionsTable struct{}

// Signature The unique signature for the migration.
func (r *M20240101000025CreateFicheConceptionsTable) Signature() string {
	return "20240101000025_create_fiche_conceptions_table"
}

// Up Run the migrations.
func (r *M20240101000025CreateFicheConceptionsTable) Up() error {
	return facades.Schema().Create("fiche_conceptions", func(table schema.Blueprint) {
		table.ID("id")
		table.String("reference", 100)
		table.Unique("reference")
		table.UnsignedBigInteger("product_variant_id").Nullable()
		table.String("title", 255)
		table.Text("description").Nullable()
		table.UnsignedBigInteger("requested_by")
		table.Enum("status", []any{"pending", "in_design", "design_done", "cancelled"}).Default("pending")
		table.Text("design_notes").Nullable()
		table.UnsignedBigInteger("validated_by").Nullable()
		table.Timestamp("validated_at").Nullable()
		table.TimestampsTz()

		table.Foreign("product_variant_id").References("id").On("product_variants")
		table.Foreign("requested_by").References("id").On("users")
		table.Foreign("validated_by").References("id").On("users")

		table.Index("product_variant_id")
		table.Index("requested_by")
		table.Index("validated_by")
		table.Index("status")
		table.Index("reference")
	})
}

// Down Reverse the migrations.
func (r *M20240101000025CreateFicheConceptionsTable) Down() error {
	return facades.Schema().DropIfExists("fiche_conceptions")
}
