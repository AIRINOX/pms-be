package migrations

import (
	"github.com/goravel/framework/contracts/database/schema"
	"github.com/goravel/framework/facades"
)

type M20240101000010CreateArticleAttributesTable struct{}

// Signature The unique signature for the migration.
func (r *M20240101000010CreateArticleAttributesTable) Signature() string {
	return "20240101000010_create_article_attributes_table"
}

// Up Run the migrations.
func (r *M20240101000010CreateArticleAttributesTable) Up() error {
	return facades.Schema().Create("article_attributes", func(table schema.Blueprint) {
		table.ID("id")
		table.UnsignedBigInteger("article_id")
		table.String("key", 100)
		table.String("title", 255)
		table.Integer("order_index").Default(0)
		table.Timestamp("created_at").UseCurrent()

		table.Foreign("article_id").References("id").On("articles")
		table.Unique("article_id", "key")
		table.Index("article_id")
		table.Index("order_index")
	})
}

// Down Reverse the migrations.
func (r *M20240101000010CreateArticleAttributesTable) Down() error {
	return facades.Schema().DropIfExists("article_attributes")
}
