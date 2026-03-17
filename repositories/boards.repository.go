package repositories

import (
	"time"

	"github.com/kasyaproject/sistem-project-management/config"
	"github.com/kasyaproject/sistem-project-management/models"
)

type BoardRepository interface {
	Create(board *models.Board) error
	FindByPublicID(publicID string) (*models.Board, error)
	Update(board *models.Board) error
	AddMember(boardID uint, userIDs []uint) error
	RemoveMember(boardID uint, userIDs []uint) error
}

type boardRepository struct {
}

func NewBoardRepository() BoardRepository {
	return &boardRepository{}
}

func (r *boardRepository) Create(board *models.Board) error {
	return config.DB.Create(board).Error
}

func (r *boardRepository) FindByPublicID(publicID string) (*models.Board, error) {
	var board models.Board

	err := config.DB.Where("public_id = ?", publicID).First(&board).Error

	return &board, err
}

func (r *boardRepository) Update(board *models.Board) error {
	return config.DB.Model(&models.Board{}).Where("public_id = ?", board.PublicID).Updates(map[string]interface{}{
		"title":       board.Title,
		"description": board.Description,
		"due_date":    board.DueDate,
	}).Error
}

func (r *boardRepository) AddMember(boardID uint, userIDs []uint) error {
	// Check if userIDs is empty
	if len(userIDs) == 0 {
		return nil
	}

	now := time.Now()                 // ambil waktu sekarang
	var members []models.BoardMembers // deklarasi variabel members

	// loop userIDs dan tambahkan data ke variabel members
	for _, userID := range userIDs {
		members = append(members, models.BoardMembers{ // tambahkan data ke variabel members
			BoardID:  int64(boardID),
			UserID:   int64(userID),
			JoinedAt: now,
		})
	}

	return config.DB.Create(&members).Error
}

func (r *boardRepository) RemoveMember(boardID uint, userIDs []uint) error {
	// Check if userIDs is empty
	if len(userIDs) == 0 {
		return nil
	}

	// Hapus data member berdasarkan boardID dan userIDs
	return config.DB.
		Where("board_internal_id = ? AND user_internal_id IN (?)", boardID, userIDs).Delete(&models.BoardMembers{}).Error
}
