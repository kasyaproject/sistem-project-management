package utils

import (
	"github.com/google/uuid"
	"github.com/kasyaproject/sistem-project-management/models"
)

func SortListByPosition(lists []models.List, order []uuid.UUID) []models.List {
	if len(order) == 0 {
		return lists
	}

	orderedList := make([]models.List, 0, len(order))

	listMap := make(map[uuid.UUID]models.List)
	for _, l := range lists {
		listMap[l.PublicID] = l
	}

	// Urutkan sesuai order nya
	for _, id := range order {
		if list, ok := listMap[id]; ok {
			orderedList = append(orderedList, list)
		}
	}

	return orderedList
}
