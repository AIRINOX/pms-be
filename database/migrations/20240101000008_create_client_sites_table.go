package migrations

import (
	"github.com/goravel/framework/contracts/database/schema"
	"github.com/goravel/framework/facades"
)

type M20240101000008CreateClientSitesTable struct{}

// Signature The unique signature for the migration.
func (r *M20240101000008CreateClientSitesTable) Signature() string {
	return "20240101000008_create_client_sites_table"
}

// Up Run the migrations.
func (r *M20240101000008CreateClientSitesTable) Up() error {
	return facades.Schema().Create("client_sites", func(table schema.Blueprint) {
		table.ID("id")
		table.UnsignedBigInteger("client_id")
		table.String("title", 255)
		table.Text("address").Nullable()
		table.String("contact_name", 255).Nullable()
		table.String("contact_phone", 50).Nullable()
		table.TimestampsTz()

		table.Foreign("client_id").References("id").On("clients")
		table.Index("client_id")
		table.Index("title")
	})
}

// Down Reverse the migrations.
func (r *M20240101000008CreateClientSitesTable) Down() error {
	return facades.Schema().DropIfExists("client_sites")
}
