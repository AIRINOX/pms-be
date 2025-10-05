package database

import (
	"github.com/goravel/framework/contracts/database/schema"
	"github.com/goravel/framework/contracts/database/seeder"

	"pms/database/migrations"
	"pms/database/seeders"
)

type Kernel struct {
}

func (kernel Kernel) Migrations() []schema.Migration {
	return []schema.Migration{
		// Base tables (no dependencies)
		&migrations.M20240101000001CreateJobsTable{},
		&migrations.M20240101000002CreateCategoriesTable{},
		&migrations.M20240101000003CreateRolesTable{},
		&migrations.M20240101000004CreateOperationsTable{},
		&migrations.M20240101000005CreateClientsTable{},
		&migrations.M20240101000006CreateStorageLocationsTable{},
		
		// Level 1 dependencies
		&migrations.M20240101000007CreateUsersTable{}, // depends on roles
		&migrations.M20240101000008CreateClientSitesTable{}, // depends on clients
		&migrations.M20240101000009CreateArticlesTable{}, // depends on categories, storage_locations
		
		// Level 2 dependencies
		&migrations.M20240101000010CreateArticleAttributesTable{}, // depends on articles
		&migrations.M20240101000011CreateArticleAttributeValuesTable{}, // depends on article_attributes
		&migrations.M20240101000012CreateArticleVariantsTable{}, // depends on articles
		&migrations.M20240101000013CreateArticleImagesTable{}, // depends on articles
		&migrations.M20240101000014CreateRecipeArticlesTable{}, // depends on articles
		
		// Level 3 dependencies
		&migrations.M20240101000015CreateOrderFabricationsTable{}, // depends on articles, article_variants, clients, client_sites, users
		&migrations.M20240101000016CreateStockMovementsTable{}, // depends on articles, article_variants, storage_locations, users
		&migrations.M20240101000017CreateStockLevelsTable{}, // depends on articles, article_variants, storage_locations
		&migrations.M20240101000018CreateRecipeArticleItemsTable{}, // depends on recipe_articles, articles
		&migrations.M20240101000019CreateRecipeVariantsTable{}, // depends on articles, article_variants
		&migrations.M20240101000020CreateRecipeVariantItemsTable{}, // depends on recipe_variants, article_variants
		&migrations.M20240101000021CreateProductionOfHistoryTable{}, // depends on order_fabrications, operations, users
		&migrations.M20240101000022CreateProductionMaterialRequirementsTable{}, // depends on order_fabrications, article_variants
		&migrations.M20240101000023CreateStockRequestsTable{}, // depends on order_fabrications, article_variants, users
		&migrations.M20240101000024CreateTechnicalDocumentsTable{}, // depends on articles, users
		&migrations.M20240101000025CreateFicheConceptionsTable{}, // depends on articles, users
	}
}

func (kernel Kernel) Seeders() []seeder.Seeder {
	return []seeder.Seeder{
		&seeders.RoleSeeder{},
		&seeders.UserSeeder{},
		&seeders.DatabaseSeeder{},
	}
}
