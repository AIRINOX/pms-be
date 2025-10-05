package migrations

import (
	"github.com/goravel/framework/contracts/database/schema"
	"github.com/goravel/framework/facades"
)

type M20240101000013CreateArticleImagesTable struct{}

// Signature The unique signature for the migration.
func (r *M20240101000013CreateArticleImagesTable) Signature() string {
	return "20240101000013_create_article_images_table"
}

// Up Run the migrations.
func (r *M20240101000013CreateArticleImagesTable) Up() error {
	return facades.Schema().Create("article_images", func(table schema.Blueprint) {
		table.ID("id")
		table.UnsignedBigInteger("article_id")
		table.String("file_path", 500)
		table.String("file_name", 255)
		table.Integer("image_index")
		table.Boolean("is_primary").Default(false)
		table.Timestamp("created_at").UseCurrent()

		table.Foreign("article_id").References("id").On("articles")
		table.Index("article_id")
		table.Index("image_index")
		table.Index("is_primary")
	})
}

// Down Reverse the migrations.
func (r *M20240101000013CreateArticleImagesTable) Down() error {
	return facades.Schema().DropIfExists("article_images")
}
