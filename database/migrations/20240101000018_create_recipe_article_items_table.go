package migrations

import (
	"github.com/goravel/framework/contracts/database/schema"
	"github.com/goravel/framework/facades"
)

type M20240101000018CreateRecipeArticleItemsTable struct{}

// Signature The unique signature for the migration.
func (r *M20240101000018CreateRecipeArticleItemsTable) Signature() string {
	return "20240101000018_create_recipe_article_items_table"
}

// Up Run the migrations.
func (r *M20240101000018CreateRecipeArticleItemsTable) Up() error {
	return facades.Schema().Create("recipe_article_items", func(table schema.Blueprint) {
		table.ID("id")
		table.UnsignedBigInteger("recipe_article_id")
		table.UnsignedBigInteger("material_article_id")
		table.Text("notes").Nullable()
		table.TimestampsTz()

		table.Foreign("recipe_article_id").References("id").On("recipe_articles")
		table.Foreign("material_article_id").References("id").On("articles")

		table.Index("recipe_article_id")
		table.Index("material_article_id")
	})
}

// Down Reverse the migrations.
func (r *M20240101000018CreateRecipeArticleItemsTable) Down() error {
	return facades.Schema().DropIfExists("recipe_article_items")
}
