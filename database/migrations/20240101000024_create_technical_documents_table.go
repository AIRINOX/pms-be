package migrations

import (
	"github.com/goravel/framework/contracts/database/schema"
	"github.com/goravel/framework/facades"
)

type M20240101000024CreateTechnicalDocumentsTable struct{}

// Signature The unique signature for the migration.
func (r *M20240101000024CreateTechnicalDocumentsTable) Signature() string {
	return "20240101000024_create_technical_documents_table"
}

// Up Run the migrations.
func (r *M20240101000024CreateTechnicalDocumentsTable) Up() error {
	return facades.Schema().Create("technical_documents", func(table schema.Blueprint) {
		table.ID("id")
		table.UnsignedBigInteger("product_variant_id").Nullable()
		table.Enum("doc_type", []any{"drawing", "spec_sheet", "quality_doc", "safety_doc"}).Nullable()
		table.UnsignedBigInteger("uploaded_by").Nullable()
		table.String("file_path", 500).Nullable()
		table.TimestampsTz()

		table.Foreign("product_variant_id").References("id").On("products")
		table.Foreign("uploaded_by").References("id").On("users")

		table.Index("product_variant_id")
		table.Index("uploaded_by")
		table.Index("doc_type")
	})
}

// Down Reverse the migrations.
func (r *M20240101000024CreateTechnicalDocumentsTable) Down() error {
	return facades.Schema().DropIfExists("technical_documents")
}
