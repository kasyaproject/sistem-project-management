package services

import (
	"errors"

	"github.com/google/uuid"
	"github.com/kasyaproject/sistem-project-management/models"
	"github.com/kasyaproject/sistem-project-management/repositories"
)

type BoardService interface {
	Create(board *models.Board) error
}

// Ambil struct dari repository
type boardService struct {
	boardRepo repositories.BoardRepository
	userRepo  repositories.UserRepository
}

func NewBoardService(
	boardRepo repositories.BoardRepository,
	userRepo repositories.UserRepository,
) BoardService {
	return &boardService{
		boardRepo,
		userRepo,
	}
}

func (s *boardService) Create(board *models.Board) error {
	// Ambil User yang sedang login
	user, err := s.userRepo.FindByPublicID(board.OwnerPublicID.String())
	if err != nil {
		return errors.New("owner not found!")
	}

	board.PublicID = uuid.New()     // generate UUID
	board.OwnerID = user.InternalID // set owner id sesuai dengan user yang sedang login

	return s.boardRepo.Create(board)
}
