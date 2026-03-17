package repositories

import (
	"time"

	"github.com/kasyaproject/sistem-project-management/config"
	"github.com/kasyaproject/sistem-project-management/models"
)

type BoardRepository interface {
	Create(board *models.Board) error
	FindByPublicID(publicID string) (*models.Board, error)
	GetAllByUserPaginate(userPublicID, filter, sort string, limit, offset int) ([]models.Board, int64, error)
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

func (r *boardRepository) GetAllByUserPaginate(userPublicID, filter, sort string, limit, offset int) ([]models.Board, int64, error) {
	var board []models.Board
	var total int64

	query := config.DB.Model(&models.Board{}).
		Where("owner_public_id = ? OR internal_id IN ("+
			"SELECT board_members.board_internal_id FROM board_members "+
			"JOIN users ON users.internal_id = board_members.user_internal_id "+
			"WHERE users.public_id = ?)", userPublicID, userPublicID)

	// Jika filter tidak kosong, tambahkan filter ke query
	if filter != "" {
		query = query.Where("title ILIKE ?", "%"+filter+"%")
	}

	// Hitung total data
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Jika sort tidak kosong, tambahkan sort ke query berdasarkan created_at
	if sort != "" {
		query = query.Order(sort)
	} else {
		query = query.Order("created_at DESC")
	}

	// Pagination data
	if err := query.Limit(limit).Offset(offset).Find(&board).Error; err != nil {
		return nil, 0, err
	}

	return board, total, nil
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
