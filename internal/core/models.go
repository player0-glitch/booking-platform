package core

import (
	"time"
	/*or uuidv7*/
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// All entities will inherit this for because of the UUIDv7 PK and FK
type Base struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	/* DeletedAt gorm.DeletedAt /*only use this for entity that will
	need it especially as an index */
}

// This allows us to create UUIDv7 for sequential ordering
func (b *Base) BeforeCreate(tx *gorm.DB) error {
	if b.ID == uuid.Nil {
		id, err := uuid.NewV7()
		if err != nil {
			return err
		}
		b.ID = id
	}
	return nil
}
