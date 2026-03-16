package repositories

import (
	"github.com/kasyaproject/sistem-project-management/config"
	"github.com/kasyaproject/sistem-project-management/models"
)

type BoardMemberRepository interface {
	GetMembers(boardPublicID string) ([]models.User, error)
}

type boardMemberRepository struct{}

func NewBoardMemberRepository() BoardMemberRepository {
	return &boardMemberRepository{}
}

func (r *boardMemberRepository) GetMembers(boardPublicID string) ([]models.User, error) {
	var user []models.User // deklarasi variabel

	// Ambil data dengan relasi antara boards dan users dengan board_members join di database
	// simpan data ke var user[]
	err := config.DB.Joins("JOIN board_members ON board_members.user_internal_id = users.internal_id").Joins("JOIN boards ON boards.internal_id = board_members.board_internal_id").Where("boards.public_id = ?", boardPublicID).Find(&user).Error

	return user, err
}
