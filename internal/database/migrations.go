package database

import (
	"fmt"
	"kepler-auth-go/internal/models"
	"log"

	"gorm.io/gorm"
)

// MigrationFunc represents a migration function
type MigrationFunc func(*gorm.DB) error

// Migration represents a database migration
type Migration struct {
	ID   string
	Up   MigrationFunc
	Down MigrationFunc
}

// GetMigrations returns all available migrations
func GetMigrations() []Migration {
	return []Migration{
		{
			ID: "001_create_initial_tables",
			Up: func(db *gorm.DB) error {
				// Create initial tables
				return db.AutoMigrate(
					&models.Organization{},
					&models.User{},
					&models.Group{},
					&models.Permission{},
					&models.AuthGroup{},
				)
			},
			Down: func(db *gorm.DB) error {
				// Drop tables in reverse order
				return db.Migrator().DropTable(
					&models.AuthGroup{},
					&models.Permission{},
					&models.Group{},
					&models.User{},
					&models.Organization{},
				)
			},
		},
		{
			ID: "002_add_indexes",
			Up: func(db *gorm.DB) error {
				// Add custom indexes for performance

				// Organizations
				if err := db.Exec("CREATE INDEX IF NOT EXISTS idx_organizations_name ON organizations(name)").Error; err != nil {
					return err
				}
				if err := db.Exec("CREATE INDEX IF NOT EXISTS idx_organizations_domain ON organizations(domain)").Error; err != nil {
					return err
				}

				// Users - composite unique index for email + organization_id
				if err := db.Exec("CREATE UNIQUE INDEX IF NOT EXISTS idx_users_email_org ON users(email, organization_id)").Error; err != nil {
					return err
				}
				if err := db.Exec("CREATE INDEX IF NOT EXISTS idx_users_organization_id ON users(organization_id)").Error; err != nil {
					return err
				}
				if err := db.Exec("CREATE INDEX IF NOT EXISTS idx_users_is_active ON users(is_active)").Error; err != nil {
					return err
				}
				if err := db.Exec("CREATE INDEX IF NOT EXISTS idx_users_is_deleted ON users(is_deleted)").Error; err != nil {
					return err
				}
				if err := db.Exec("CREATE INDEX IF NOT EXISTS idx_users_is_verified ON users(is_verified)").Error; err != nil {
					return err
				}

				// Groups - composite unique index for name + organization_id
				if err := db.Exec("CREATE UNIQUE INDEX IF NOT EXISTS idx_groups_name_org ON groups(name, organization_id)").Error; err != nil {
					return err
				}
				if err := db.Exec("CREATE INDEX IF NOT EXISTS idx_groups_organization_id ON groups(organization_id)").Error; err != nil {
					return err
				}
				if err := db.Exec("CREATE INDEX IF NOT EXISTS idx_groups_is_active ON groups(is_active)").Error; err != nil {
					return err
				}

				return nil
			},
			Down: func(db *gorm.DB) error {
				// Drop indexes

				// Organizations
				if err := db.Exec("DROP INDEX IF EXISTS idx_organizations_name").Error; err != nil {
					return err
				}
				if err := db.Exec("DROP INDEX IF EXISTS idx_organizations_domain").Error; err != nil {
					return err
				}

				// Users
				if err := db.Exec("DROP INDEX IF EXISTS idx_users_email_org").Error; err != nil {
					return err
				}
				if err := db.Exec("DROP INDEX IF EXISTS idx_users_organization_id").Error; err != nil {
					return err
				}
				if err := db.Exec("DROP INDEX IF EXISTS idx_users_is_active").Error; err != nil {
					return err
				}
				if err := db.Exec("DROP INDEX IF EXISTS idx_users_is_deleted").Error; err != nil {
					return err
				}
				if err := db.Exec("DROP INDEX IF EXISTS idx_users_is_verified").Error; err != nil {
					return err
				}

				// Groups
				if err := db.Exec("DROP INDEX IF EXISTS idx_groups_name_org").Error; err != nil {
					return err
				}
				if err := db.Exec("DROP INDEX IF EXISTS idx_groups_organization_id").Error; err != nil {
					return err
				}
				if err := db.Exec("DROP INDEX IF EXISTS idx_groups_is_active").Error; err != nil {
					return err
				}

				return nil
			},
		},
	}
}

// MigrationRecord tracks which migrations have been applied
type MigrationRecord struct {
	ID        uint   `gorm:"primaryKey"`
	Migration string `gorm:"uniqueIndex"`
	AppliedAt int64  `gorm:"autoCreateTime"`
}

// RunMigrations executes all pending migrations
func RunMigrations(force bool) error {
	// Create migration tracking table
	if err := DB.AutoMigrate(&MigrationRecord{}); err != nil {
		return fmt.Errorf("failed to create migration table: %w", err)
	}

	migrations := GetMigrations()

	for _, migration := range migrations {
		// Check if migration already applied
		var record MigrationRecord
		err := DB.Where("migration = ?", migration.ID).First(&record).Error

		if err == nil && !force {
			log.Printf("Migration %s already applied, skipping", migration.ID)
			continue
		}

		if err != nil && err != gorm.ErrRecordNotFound {
			return fmt.Errorf("failed to check migration status: %w", err)
		}

		log.Printf("Running migration: %s", migration.ID)

		// Start transaction
		tx := DB.Begin()
		if tx.Error != nil {
			return fmt.Errorf("failed to start transaction: %w", tx.Error)
		}

		// Run migration
		if err := migration.Up(tx); err != nil {
			tx.Rollback()
			return fmt.Errorf("migration %s failed: %w", migration.ID, err)
		}

		// Record migration
		migrationRecord := MigrationRecord{Migration: migration.ID}
		if err := tx.Create(&migrationRecord).Error; err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to record migration: %w", err)
		}

		// Commit transaction
		if err := tx.Commit().Error; err != nil {
			return fmt.Errorf("failed to commit migration: %w", err)
		}

		log.Printf("Migration %s completed successfully", migration.ID)
	}

	return nil
}

// RollbackMigration rolls back the last migration
func RollbackMigration() error {
	// Get the last applied migration
	var record MigrationRecord
	err := DB.Order("applied_at DESC").First(&record).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return fmt.Errorf("no migrations to rollback")
		}
		return fmt.Errorf("failed to get last migration: %w", err)
	}

	// Find the migration
	migrations := GetMigrations()
	var targetMigration *Migration
	for _, migration := range migrations {
		if migration.ID == record.Migration {
			targetMigration = &migration
			break
		}
	}

	if targetMigration == nil {
		return fmt.Errorf("migration %s not found", record.Migration)
	}

	log.Printf("Rolling back migration: %s", targetMigration.ID)

	// Start transaction
	tx := DB.Begin()
	if tx.Error != nil {
		return fmt.Errorf("failed to start transaction: %w", tx.Error)
	}

	// Run rollback
	if err := targetMigration.Down(tx); err != nil {
		tx.Rollback()
		return fmt.Errorf("rollback %s failed: %w", targetMigration.ID, err)
	}

	// Remove migration record
	if err := tx.Delete(&record).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to remove migration record: %w", err)
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("failed to commit rollback: %w", err)
	}

	log.Printf("Migration %s rolled back successfully", targetMigration.ID)
	return nil
}

// GetMigrationStatus returns the status of all migrations
func GetMigrationStatus() ([]map[string]interface{}, error) {
	migrations := GetMigrations()
	var records []MigrationRecord

	if err := DB.Find(&records).Error; err != nil {
		return nil, fmt.Errorf("failed to get migration records: %w", err)
	}

	// Create a map for quick lookup
	appliedMap := make(map[string]MigrationRecord)
	for _, record := range records {
		appliedMap[record.Migration] = record
	}

	var status []map[string]interface{}
	for _, migration := range migrations {
		record, applied := appliedMap[migration.ID]
		status = append(status, map[string]interface{}{
			"id":      migration.ID,
			"applied": applied,
			"applied_at": func() interface{} {
				if applied {
					return record.AppliedAt
				}
				return nil
			}(),
		})
	}

	return status, nil
}
