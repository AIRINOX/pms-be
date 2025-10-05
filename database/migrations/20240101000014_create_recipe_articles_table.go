package migrations

import (
	"github.com/goravel/framework/contracts/database/schema"
	"github.com/goravel/framework/facades"
)

type M20240101000014CreateRecipeArticlesTable struct{}

// Signature The unique signature for the migration.
func (r *M20240101000014CreateRecipeArticlesTable) Signature() string {
	return "20240101000014_create_recipe_articles_table"
}

// Up Run the migrations.
func (r *M20240101000014CreateRecipeArticlesTable) Up() error {
	return facades.Schema().Create("recipe_articles", func(table schema.Blueprint) {
		table.ID("id")
		table.String("title", 255)
		table.Text("description").Nullable()
		table.UnsignedBigInteger("article_id")
		table.TimestampsTz()

		table.Foreign("article_id").References("id").On("articles")
		table.Index("article_id")
		table.Index("title")
	})
}

// Down Reverse the migrations.
func (r *M20240101000014CreateRecipeArticlesTable) Down() error {
	return facades.Schema().DropIfExists("recipe_articles")
}
