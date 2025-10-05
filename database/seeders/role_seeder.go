package seeders

import (
	"github.com/goravel/framework/facades"

	"pms/app/models"
)

type RoleSeeder struct {
}

// Signature The unique signature for the seeder.
func (s *RoleSeeder) Signature() string {
	return "RoleSeeder"
}

// Run executes the seeder.
func (s *RoleSeeder) Run() error {
	// Check if roles already exist
	count, err := facades.Orm().Query().Model(&models.Role{}).Count()
	if err != nil {
		facades.Log().Error("Failed to count roles")
		return err
	}

	if count > 0 {
		facades.Log().Info("Roles already exist, skipping seeder")
		return nil
	}

	// Create roles for Production Management system
	roles := []models.Role{
		{
			Key:        "admin",
			Title:      "Admin",
			OrderIndex: 1,
		},
		{
			Key:        "commercial",
			Title:      "Commercial",
			OrderIndex: 2,
		},
		{
			Key:        "ingenieur_methodes",
			Title:      "Ingénieur Méthodes",
			OrderIndex: 3,
		},
		{
			Key:        "magasinier",
			Title:      "Magasinier",
			OrderIndex: 4,
		},
		{
			Key:        "achat",
			Title:      "Achat",
			OrderIndex: 5,
		},
		{
			Key:        "operateur_decoupe",
			Title:      "Opérateur de Découpe",
			OrderIndex: 6,
		},
		{
			Key:        "operateur_pliage",
			Title:      "Opérateur de Pliage",
			OrderIndex: 7,
		},
		{
			Key:        "operateur_assemblage",
			Title:      "Opérateur d'Assemblage",
			OrderIndex: 8,
		},
		{
			Key:        "operateur_finition",
			Title:      "Opérateur de Finition",
			OrderIndex: 9,
		},
	}

	for _, role := range roles {
		if err := facades.Orm().Query().Create(&role); err != nil {
			facades.Log().Error("Failed to create role: " + role.Key)
			return err
		}
		facades.Log().Info("Created role: " + role.Key)
	}

	return nil
}
