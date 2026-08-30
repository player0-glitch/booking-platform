// database/migrator.go
package database

import (
	"fmt"

	"gorm.io/gorm"
)

// I'm not sure i understand this yet
type Migration func(*gorm.DB) error

type Migrator struct {
	dbContext  *gorm.DB
	migrations []Migration
}

// It's a constructor, you can tell by the return type of the pointer
func NewMigrator(db *gorm.DB, migrations ...Migration) *Migrator {
	//variadic so accepts arbitrary migrations
	return &Migrator{dbContext: db, migrations: migrations}
}

func (migrator *Migrator) Run() error {
	for _, migration := range migrator.migrations {
		err := migration(migrator.dbContext)
		if err != nil {
			return fmt.Errorf("Migrations failed %w", err)
		}
	}
	return nil
}
