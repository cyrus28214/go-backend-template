package user

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BaseModel struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (b *BaseModel) BeforeCreate(tx *gorm.DB) (err error) {
	b.ID, err = uuid.NewV7()
	return err
}

type User struct {
	BaseModel
	Email    string `gorm:"uniqueIndex;not null;size:255"` // RFC 5321, longest email address is 254 characters
	Password string `gorm:"not null"`
}
