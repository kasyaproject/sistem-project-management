package utils

import (
	"sort"

	"github.com/google/uuid"
	"github.com/kasyaproject/sistem-project-management/models"
)

func SortCardByPosition(cards []models.Card, order []uuid.UUID) []models.Card {
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
