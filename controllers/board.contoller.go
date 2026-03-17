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

func (c *BoardController) UpdateBoard(ctx *fiber.Ctx) error {
	publicID := ctx.Params("id")
	board := new(models.Board)

	// Check input data
	if err := ctx.BodyParser(board); err != nil {
		return utils.BadRequest(ctx, "Gagal Parsing data", err.Error())
	}

	// Validasi ID
	if _, err := uuid.Parse(publicID); err != nil {
		return utils.BadRequest(ctx, "ID tidak valid", err.Error())
	}

	// Ambil data board
	existingBoard, err := c.service.GetByPublicID(publicID)
	if err != nil {
		return utils.NotFound(ctx, "Board not found", err.Error())
	}

	board.InternalID = existingBoard.InternalID // set internal id sesuai dengan board yang sudah ada
	board.PublicID = existingBoard.PublicID     // set public id sesuai dengan board yang sudah ada

	// Update board
	if err := c.service.Update(board); err != nil {
		return utils.BadRequest(ctx, "Gagal Update board", err.Error())
	}

	return utils.Success(ctx, "Update board success", board)
}

func (c *BoardController) AddBoardMember(ctx *fiber.Ctx) error {
	publicID := ctx.Params("id")
	var userIDs []string

	if err := ctx.BodyParser(&userIDs); err != nil {
		return utils.BadRequest(ctx, "Gagal Parsing data", err.Error())
	}

	if err := c.service.AddMembers(publicID, userIDs); err != nil {
		return utils.BadRequest(ctx, "Gagal Add members", err.Error())
	}

	return utils.Success(ctx, "Add members success", nil)
}

func (c *BoardController) RemoveBoardMember(ctx *fiber.Ctx) error {
	publicID := ctx.Params("id")
	var userIDs []string

	if err := ctx.BodyParser(&userIDs); err != nil {
		return utils.BadRequest(ctx, "Gagal Parsing data", err.Error())
	}

	if err := c.service.RemoveMembers(publicID, userIDs); err != nil {
		return utils.BadRequest(ctx, "Gagal Remove members", err.Error())
	}

	return utils.Success(ctx, "Remove members success", nil)
}
