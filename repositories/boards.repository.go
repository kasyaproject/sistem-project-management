package repositories

import (
	"github.com/kasyaproject/sistem-project-management/config"
	"github.com/kasyaproject/sistem-project-management/models"
)

type BoardRepository interface {
	Create(board *models.Board) error
}

type boardRepository struct {
}

func NewBoardRepository() BoardRepository {
	return &boardRepository{}
}

func (r boardRepository) Create(board *models.Board) error {
	return config.DB.Create(board).Error
}
