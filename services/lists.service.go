package services

import (
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/kasyaproject/sistem-project-management/config"
	"github.com/kasyaproject/sistem-project-management/models"
	"github.com/kasyaproject/sistem-project-management/models/types"
	"github.com/kasyaproject/sistem-project-management/repositories"
	"github.com/kasyaproject/sistem-project-management/utils"
	"gorm.io/gorm"
)

type ListService interface {
	GetByBoardID(boardPublicID string) (*ListWithOrder, error)
	GetByID(id uint) (*models.List, error)
	GetByPublicID(publicID string) (*models.List, error)

	Create(list *models.List) error
	Update(list *models.List) error
	Delete(id uint) error
	UpdatePosition(boardPublicID string, position []uuid.UUID) error
}

type ListWithOrder struct {
	Positions []uuid.UUID
	Lists     []models.List
}

type listService struct {
	listRepo         repositories.ListRepository
	boardRepo        repositories.BoardRepository
	listPositionRepo repositories.ListPositionsRepository
}

func NewListService(
	listRepo repositories.ListRepository,
	boardRepo repositories.BoardRepository,
	listPositionRepo repositories.ListPositionsRepository,
) ListService {
	return &listService{
		listRepo,
		boardRepo,
		listPositionRepo,
	}
}

func (s *listService) GetByBoardID(boardPublicID string) (*ListWithOrder, error) {
	// Cek apakah board ada
	_, err := s.boardRepo.FindByPublicID(boardPublicID)
	if err != nil {
		return nil, errors.New("board not found")
	}

	// Ambil Positions list dari board
	position, err := s.listPositionRepo.GetListOrder(boardPublicID)
	// Jika position kosong
	if len(position) == 0 {
		return nil, errors.New("List position empty")
	}

	if err != nil {
		return nil, errors.New("failed to get list order : " + err.Error())
	}

	// Ambil list berdasarkan boardID
	lists, err := s.listRepo.FindByBoardID(boardPublicID)
	if err != nil {
		return nil, errors.New("failed to get list" + err.Error())
	}

	// Sorting berdasarkan position
	orderedList := utils.SortListByPosition(lists, position)

	return &ListWithOrder{
		Positions: position,
		Lists:     orderedList,
	}, nil
}

func (s *listService) GetByID(id uint) (*models.List, error) {
	return s.listRepo.FindByID(id)
}

func (s *listService) GetByPublicID(publicID string) (*models.List, error) {
	return s.listRepo.FindByPublicID(publicID)
}

func (s *listService) Create(list *models.List) error {
	// Validasi apakah board ada atau tidak
	board, err := s.boardRepo.FindByPublicID(list.BoardPublicID.String())
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("board not found!!!")
		}
		return fmt.Errorf("failed to get board: %w", err)
	}

	list.BoardInternalID = board.InternalID // Set board internalID ke list

	// Validasi apakah publicID ada atau tidak dan generate jika tidak ada
	if list.PublicID == uuid.Nil {
		list.PublicID = uuid.New()
	}

	// Mulai transaction
	tx := config.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Simpan list baru
	if err := tx.Create(list).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to create list : %w", err)
	}

	// update position list
	var position models.ListPosition
	res := tx.Where("board_internal_id = ?", board.InternalID).First(&position)

	// Cek apakah sudah ada atau belum
	if errors.Is(res.Error, gorm.ErrRecordNotFound) {
		// Jika belum ada maka buat baru
		position = models.ListPosition{
			PublicID:  uuid.New(),
			BoardID:   board.InternalID,
			ListOrder: types.UUIDArray{list.PublicID},
		}

		if err := tx.Create(&position).Error; err != nil {
			tx.Rollback()
			return fmt.Errorf("(1)failed to create list position: %w", err)
		}
	} else if res.Error != nil {
		tx.Rollback()
		return fmt.Errorf("(2)failed to create list position: %w", res.Error)
	} else {
		// Jika sudah ada maka update position
		position.ListOrder = append(position.ListOrder, list.PublicID)
		// update ke DB
		if err := tx.Model(&position).Update("list_order", position.ListOrder).Error; err != nil {
			tx.Rollback()
			return fmt.Errorf("(3)failed to update list position: %w", err)
		}
	}

	// commit transaction
	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (s *listService) Update(list *models.List) error {
	return s.listRepo.Update(list)
}

func (s *listService) Delete(id uint) error {
	return s.listRepo.Delete(id)
}

func (s *listService) UpdatePosition(boardPublicID string, position []uuid.UUID) error {
	// Verifikasi apakah board ada
	board, err := s.boardRepo.FindByPublicID(boardPublicID)
	if err != nil {
		return errors.New("board not found")
	}

	// Get list order position sekarang
	listPosition, err := s.listPositionRepo.GetByBoard(board.PublicID.String())
	if err != nil {
		return errors.New("failed to get list position" + err.Error())
	}

	// Update list order position
	listPosition.ListOrder = position
	// Update list order position baru
	return s.listPositionRepo.UpdateListOrder(listPosition)
}
