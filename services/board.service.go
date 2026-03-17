package services

import (
	"errors"

	"github.com/google/uuid"
	"github.com/kasyaproject/sistem-project-management/models"
	"github.com/kasyaproject/sistem-project-management/repositories"
)

type BoardService interface {
	Create(board *models.Board) error
	GetByPublicID(publicID string) (*models.Board, error)
	Update(board *models.Board) error
	AddMembers(boardPublicID string, userPublicIDs []string) error
	RemoveMembers(boardPublicID string, userPublicIDs []string) error
}

// Ambil struct dari repository
type boardService struct {
	boardRepo       repositories.BoardRepository
	userRepo        repositories.UserRepository
	boardMemberRepo repositories.BoardMemberRepository
}

func NewBoardService(
	boardRepo repositories.BoardRepository,
	userRepo repositories.UserRepository,
	boardMemberRepo repositories.BoardMemberRepository,
) BoardService {
	return &boardService{
		boardRepo,
		userRepo,
		boardMemberRepo,
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

func (s *boardService) GetByPublicID(publicID string) (*models.Board, error) {
	return s.boardRepo.FindByPublicID(publicID)
}

func (s *boardService) Update(board *models.Board) error {
	return s.boardRepo.Update(board)
}

func (s *boardService) AddMembers(boardPublicID string, userPublicIDs []string) error {
	// Ambil data board
	board, err := s.boardRepo.FindByPublicID(boardPublicID)
	if err != nil {
		return errors.New("board not found!")
	}

	// Ambil data internalID user dari table users
	var userInternalIDs []uint                   // deklarasi variabel
	for _, userPublicID := range userPublicIDs { // loop userPublicIDs
		user, err := s.userRepo.FindByPublicID(userPublicID) // Cari user berdasarkan publicID

		// Jika user tidak ditemukan
		if err != nil {
			return errors.New("user not found: " + userPublicID)
		}

		// Tambahkan internalID user ke variabel userInternalID
		userInternalIDs = append(userInternalIDs, uint(user.InternalID))
	}

	// Ambil data member di board
	existingMembers, err := s.boardMemberRepo.GetMembers(string(board.PublicID.String()))
	if err != nil {
		return errors.New("board not found!")
	}

	// Buat map untuk menyimpan internalID member
	memberMap := make(map[uint]bool)
	for _, member := range existingMembers {
		memberMap[uint(member.InternalID)] = true // memberMap[internalID] = true
	}

	// Simpan member yang ingin ditambahkan, kedalam variable newMemberIDs
	var newMemberIDs []uint
	for _, userID := range userInternalIDs {
		if !memberMap[userID] { // Jika memberMap[internalID] tidak ada
			newMemberIDs = append(newMemberIDs, userID) // Tambahkan memberID ke variabel newMemberIDs
		}
	}

	// Jika variabel newMemberIDs masih kosong
	if len(newMemberIDs) == 0 {
		return nil
	}

	return s.boardRepo.AddMember(uint(board.InternalID), newMemberIDs)
}

func (s *boardService) RemoveMembers(boardPublicID string, userPublicIDs []string) error {
	// Ambil data board
	board, err := s.boardRepo.FindByPublicID(boardPublicID)
	if err != nil {
		return errors.New("board not found!")
	}

	// Ambil data internalID user dari table users
	var userInternalIDs []uint                   // deklarasi variabel
	for _, userPublicID := range userPublicIDs { // loop userPublicIDs
		user, err := s.userRepo.FindByPublicID(userPublicID) // Cari user berdasarkan publicID

		// Jika user tidak ditemukan
		if err != nil {
			return errors.New("user not found: " + userPublicID)
		}

		// Tambahkan internalID user ke variabel userInternalID
		userInternalIDs = append(userInternalIDs, uint(user.InternalID))
	}

	// Ambil data member di board
	existingMembers, err := s.boardMemberRepo.GetMembers(string(board.PublicID.String()))
	if err != nil {
		return errors.New("board member not found!")
	}

	// Buat map untuk menyimpan internalID member
	memberMap := make(map[uint]bool)
	for _, member := range existingMembers {
		memberMap[uint(member.InternalID)] = true // memberMap[internalID] = true
	}

	// Simpan member yang ingin dihapus, kedalam variable membersToRemove
	var membersToRemove []uint
	for _, userID := range userInternalIDs {
		if memberMap[userID] {
			membersToRemove = append(membersToRemove, userID)
		}
	}

	return s.boardRepo.RemoveMember(uint(board.InternalID), membersToRemove)
}
