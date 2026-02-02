package models

import (
	"github.com/google/uuid"
	"github.com/kasyaproject/sistem-project-management/models/types"
)

type Card_Position struct {
	InternalID int64     `json:"internal_id" gorm:"primaryKey;autoIncrement"`
	PublicID   uuid.UUID `json:"public_id" gorm:"type:uuid;not null"`
	ListID     int64     `json:"list_internal_id" gorm:"column:list_internal_id:not null"`
	// Menggunakan custom type UUIDArray untuk menyimpan urutan list dalam bentuk string UUID
	CardOrder types.UUIDArray `json:"card_order" gorm:"type:uuid[]"`
}
