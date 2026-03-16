package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"github.com/kasyaproject/sistem-project-management/models"
	"github.com/kasyaproject/sistem-project-management/services"
	"github.com/kasyaproject/sistem-project-management/utils"
)

type BoardController struct {
	service services.BoardService
}

func NewBoardController(s services.BoardService) *BoardController {
	return &BoardController{service: s}
}

func (c *BoardController) CreateBoard(ctx *fiber.Ctx) error {
	var userID uuid.UUID
	board := new(models.Board)

	// Ambil data user yang sedang login
	user := ctx.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)

	// Check input data
	if err := ctx.BodyParser(board); err != nil {
		return utils.BadRequest(ctx, "Gagal Parsing data", err.Error())
	}

	// Ambil data user yang sedang login
	userID, err := uuid.Parse(claims["pub_id"].(string))
	if err != nil {
		return utils.BadRequest(ctx, "Gagal Parsing data jwt", err.Error())
	}

	// Set owner id sesuai dengan user yang sedang login
	board.OwnerPublicID = userID

	// Create board
	if err := c.service.Create(board); err != nil {
		return utils.BadRequest(ctx, "Gagal Create board", err.Error())
	}

	return utils.Success(ctx, "Create board success", board)
}
