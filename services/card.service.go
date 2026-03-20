package services

import (
	"errors"
	"fmt"
	"sort"
	"time"

	"github.com/google/uuid"
	"github.com/kasyaproject/sistem-project-management/config"
	"github.com/kasyaproject/sistem-project-management/models"
	"github.com/kasyaproject/sistem-project-management/models/types"
	"github.com/kasyaproject/sistem-project-management/repositories"
	"github.com/kasyaproject/sistem-project-management/utils"
	"gorm.io/gorm"
)

type CardService interface {
	Create(card *models.Card, listPublicID string) error
	Update(card *models.Card, listPublicID string) error
	Delete(id uint) error

	GetByListID(listPublicID string) ([]models.Card, error)
	GetByID(id uint) (*models.Card, error)
	GetByPublicID(publicID string) (*models.Card, error)
}

type cardService struct {
	cardRepo repositories.CardRepository
	listRepo repositories.ListRepository
	userRepo repositories.UserRepository
}

func NewCardService(
	cardRepo repositories.CardRepository,
	listRepo repositories.ListRepository,
	userRepo repositories.UserRepository,
) CardService {
	return &cardService{
		cardRepo,
		listRepo,
		userRepo,
	}
}

func (s *cardService) Create(card *models.Card, listPublicID string) error {
	// Ambil list dari listPublicID
	list, err := s.listRepo.FindByPublicID(listPublicID)
	if err != nil {
		return fmt.Errorf("list not found : %w", err)
	}

	card.ListID = list.InternalID  // Set list internalID ke card
	card.CreatedAt = time.Now()    // Set created_at ke card
	if card.PublicID == uuid.Nil { // Generate public_id jika belum ada
		card.PublicID = uuid.New()
	}

	// Mulai transaction
	tx := config.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
	}()

	// Simpan card baru
	if err := tx.Create(card).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("1. failed to create card : %w", err)
	}

	// Update atau buat card_position
	var position models.CardPosition
	if err := tx.Model(&models.CardPosition{}).
		Where("list_internal_id = ?", list.InternalID).
		First(&position).Error; err != nil {

		// Kondisi jika card position belum ada
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// Buat card position jika belum ada
			position = models.CardPosition{
				PublicID:  uuid.New(),
				ListID:    list.InternalID,
				CardOrder: types.UUIDArray{card.PublicID},
			}

			// Simpan card position baru
			if err := tx.Model(&position).Error; err != nil {
				tx.Rollback()
				return fmt.Errorf("failed to create card position : %w", err)
			}
		} else {
			tx.Rollback()
			return fmt.Errorf("failed to get card position : %w", err)
		}
	} else {
		// Kondisi jika card position sudah ada
		// Tambah card baru ke urutan
		position.CardOrder = append(position.CardOrder, card.PublicID)

		// Update card position
		if err := tx.Model(&models.CardPosition{}).
			Where("internal_id = ?", position.InternalID).
			Update("card_order", position.CardOrder).Error; err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to update card position : %w", err)
		}
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("failed to commit transaction : %w", err)
	}

	return nil
}

func (s *cardService) Update(card *models.Card, listPublicID string) error {
	// Ambil data card
	existingCard, err := s.cardRepo.FindByPublicID(card.PublicID.String())
	if err != nil {
		return fmt.Errorf("card not found : %w", err)
	}

	// Ambil data list
	newList, err := s.listRepo.FindByPublicID(listPublicID)
	if err != nil {
		return fmt.Errorf("list not found : %w", err)
	}

	// Mulai Transaction
	tx := config.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
	}()

	// jika card pindah list -> hapus card dari posisi list lama & tambah list baru
	if existingCard.ListID != newList.InternalID { // pindah posisi list
		// hapus dari posisi lama
		var oldPos models.CardPosition
		if err := tx.Where("list_internal_id = ?", existingCard.ListID).First(&oldPos).Error; err != nil {
			// Filter card position di list lama sebelumnya tanpa card yang dipindahkan
			// karena disini ingin menghapus card yang dipindahkan dari list lama
			filtered := make(types.UUIDArray, 0, len(oldPos.CardOrder))
			for _, id := range oldPos.CardOrder {
				if id != existingCard.PublicID {
					filtered = append(filtered, id)
				}
			}

			// Update card position
			if err := tx.Model(&models.CardPosition{}).
				Where("internal_id = ?", oldPos.InternalID).
				Update("card_order", types.UUIDArray(filtered)).Error; err != nil {
				tx.Rollback()
				return fmt.Errorf("failed to update card position : %w", err)
			}
		} else if !errors.Is(err, gorm.ErrRecordNotFound) {
			tx.Rollback()
			return fmt.Errorf("failed to get old card position : %w", err)
		}

		// Tambah card ke list baru
		var newPos models.CardPosition
		// Ambil data card position di list baru lalu simpan di var newPos
		res := tx.Where("list_internal_id = ?", newList.InternalID).First(&newPos)

		// kondisi jika di list baru belum ada card position
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			// Buat card position di list yang baru
			newPos = models.CardPosition{
				PublicID:  uuid.New(),
				ListID:    newList.InternalID,
				CardOrder: types.UUIDArray{existingCard.PublicID},
			}

			// Simpan card position baru di list baru
			if err := tx.Model(&newPos).Error; err != nil {
				tx.Rollback()
				return fmt.Errorf("failed to create card position for new list : %w", err)
			}
		} else if res.Error == nil {
			// kondisi jika di list baru sudah ada card position
			// maka appen card baru ke card order di list baru
			updateOrder := append(newPos.CardOrder, existingCard.PublicID)

			if err := tx.Model(&models.CardPosition{}).
				Where("internal_id = ?", newPos.InternalID).
				Update("card_order", types.UUIDArray(updateOrder)).Error; err != nil {
				tx.Rollback()
				return fmt.Errorf("failed to update card position for new list : %w", err)
			}
		} else {
			// Kondisi jika error
			tx.Rollback()
			return fmt.Errorf("failed to get new card position : %w", err)
		}
	}

	// Update data card
	card.InternalID = existingCard.InternalID
	card.PublicID = existingCard.PublicID
	card.ListID = existingCard.ListID

	if err := tx.Save(card).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to update card : %w", err)
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("Transaction commit failed : %w", err)
	}

	return nil
}

func (s *cardService) Delete(id uint) error {
	return s.cardRepo.Delete(id)
}

func (s *cardService) GetByListID(listPublicID string) ([]models.Card, error) {
	// Ambil data list
	list, err := s.listRepo.FindByPublicID(listPublicID)
	if err != nil {
		return nil, fmt.Errorf("list not found : %w", err)
	}

	// Ambil card position
	position, err := s.cardRepo.FindCardPositionByListID(list.InternalID)
	if err != nil {
		return nil, fmt.Errorf("failed to get card position : %w", err)
	}

	// Ambil data card berdasarkan list
	cards, err := s.cardRepo.FindByListID(listPublicID)
	if err != nil {
		return nil, fmt.Errorf("failed to get card : %w", err)
	}

	// Sort card berdasarkan position
	if position != nil && len(position.CardOrder) > 0 {
		cards = utils.SortCardByPosition(cards, position.CardOrder)
	}

	return cards, nil
}

func sortCardByPosition(cards []models.Card, order []uuid.UUID) []models.Card {
	// Buat map untuk pencarian cepat
	orderMap := make(map[uuid.UUID]int)
	for i, id := range order {
		orderMap[id] = i
	}

	defaultIndex := len(order)

	// Sort card berdasarkan position dan created_at
	sort.SliceStable(cards, func(i, j int) bool {
		indxI, okI := orderMap[cards[i].PublicID]
		if !okI {
			indxI = defaultIndex
		}

		indxJ, okJ := orderMap[cards[j].PublicID]
		if !okJ {
			indxJ = defaultIndex
		}

		if indxI == indxJ {
			return cards[i].CreatedAt.Before(cards[j].CreatedAt)
		}

		return indxI < indxJ
	})

	return cards
}

func (s *cardService) GetByID(id uint) (*models.Card, error) {
	return s.cardRepo.FindByID(id)
}

func (s *cardService) GetByPublicID(publicID string) (*models.Card, error) {
	return s.cardRepo.FindByPublicID(publicID)
}
