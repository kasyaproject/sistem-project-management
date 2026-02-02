package models

import (
	"time"

	"github.com/google/uuid"
)

type Board struct {
	InternalId    int64      `json:"internal_id" db:"internal_id" gorm:"primaryKey;autoIncrement"`
	PublicId      uuid.UUID  `json:"public_id" db:"public_id"`
	Title         string     `json:"title" db:"title"`
	Description   string     `json:"description" db:"description"`
	Owner         int64      `json:"owner_internal_id" db:"owner_internal_id" gorm:"column: owner_internal_id"`
	OwnerPublicId uuid.UUID  `json:"owner_public_id" db:"owner_public_id"`
	CreatedAt     time.Time  `json:"created_at" db:"created_at"`
	Duedate       *time.Time `json:"due_date,omitempty" db:"due_date"`
}
