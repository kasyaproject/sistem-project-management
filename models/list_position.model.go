package models

import (
	"github.com/google/uuid"
	"github.com/kasyaproject/sistem-project-management/models/types"
)

type List_Position struct {
	InternalID int64     `json:"internal_id" db:"internal_id" gorm:"primaryKey;autoIncrement"`
	PublicID   uuid.UUID `json:"public_id" db:"public_id"`
	BoardID    int64     `json:"board_internal_id" db:"board_internal_id" gorm:"column: board_internal_id"`
	// Menggunakan custom type UUIDArray untuk menyimpan urutan list dalam bentuk string UUID
	ListOrder types.UUIDArray `json:"list_order"`
}
