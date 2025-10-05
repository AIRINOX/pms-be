package migrations

import (
	"github.com/goravel/framework/contracts/database/schema"
	"github.com/goravel/framework/facades"
)

type M20240101000011CreateArticleAttributeValuesTable struct{}

// Signature The unique signature for the migration.
func (r *M20240101000011CreateArticleAttributeValuesTable) Signature() string {
	return "20240101000011_create_article_attribute_values_table"
}

// Up Run the migrations.
func (r *M20240101000011CreateArticleAttributeValuesTable) Up() error {
	return facades.Schema().Create("article_attribute_values", func(table schema.Blueprint) {
		table.ID("id")
		table.UnsignedBigInteger("attribute_id")
		table.String("value", 255)
		table.Integer("order_index").Default(0)
		table.Boolean("is_active").Default(true)
		table.Timestamp("created_at").UseCurrent()

		table.Foreign("attribute_id").References("id").On("article_attributes")
		table.Index("attribute_id")
		table.Index("order_index")
		table.Index("is_active")
	})
}

// Down Reverse the migrations.
func (r *M20240101000011CreateArticleAttributeValuesTable) Down() error {
	return facades.Schema().DropIfExists("article_attribute_values")
}
