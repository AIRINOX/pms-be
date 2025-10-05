package migrations

import (
	"github.com/goravel/framework/contracts/database/schema"
	"github.com/goravel/framework/facades"
)

type M20240101000007CreateUsersTable struct{}

// Signature The unique signature for the migration.
func (r *M20240101000007CreateUsersTable) Signature() string {
	return "20240101000007_create_users_table"
}

// Up Run the migrations.
func (r *M20240101000007CreateUsersTable) Up() error {
	return facades.Schema().Create("users", func(table schema.Blueprint) {
		table.ID("id")
		table.String("username", 100)
		table.Unique("username")
		table.String("password", 255)
		table.String("name", 255).Nullable()
		table.String("email", 255).Nullable()
		table.String("phone", 50).Nullable()
		table.UnsignedBigInteger("role_id")
		table.Boolean("is_active").Default(true)
		table.TimestampsTz()

		table.Foreign("role_id").References("id").On("roles")
		table.Index("username")
		table.Index("role_id")
		table.Index("is_active")
	})
}

// Down Reverse the migrations.
func (r *M20240101000007CreateUsersTable) Down() error {
	return facades.Schema().DropIfExists("users")
}
