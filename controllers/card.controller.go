package controllers

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/kasyaproject/sistem-project-management/models"
	"github.com/kasyaproject/sistem-project-management/services"
	"github.com/kasyaproject/sistem-project-management/utils"
)

type CardController struct {
	service services.CardService
}

func NewCardController(s services.CardService) *CardController {
	return &CardController{service: s}
}

func (c *CardController) CreateCard(ctx *fiber.Ctx) error {
	// Struct request
	type CreateCardRequest struct {
		ListPublicID string    `json:"list_id"`
		Title        string    `json:"title"`
		Description  string    `json:"description"`
		DueDate      time.Time `json:"due_date"`
		Position     int       `json:"position"`
	}

	// Parse request
	var req CreateCardRequest
	if err := ctx.BodyParser(&req); err != nil {
		return utils.BadRequest(ctx, "Gagal parsing data", err.Error())
	}

	// Tampung data request dalam struct card yang sesuai dengan model card
	card := &models.Card{
		Title:       req.Title,
		Description: req.Description,
		DueDate:     &req.DueDate,
		Position:    req.Position,
	}

	// Buat card
	if err := c.service.Create(card, req.ListPublicID); err != nil {
		return utils.InternalServerError(ctx, "Gagal membuat card", err.Error())
	}

	return utils.Success(ctx, "Card berhasil dibuat", card)
}

func (c *CardController) UpdateCard(ctx *fiber.Ctx) error {
	publicID := ctx.Params("id")

	// Struct request
	type UpdateCardRequest struct {
		ListPublicID string    `json:"list_id"`
		Title        string    `json:"title"`
		Description  string    `json:"description"`
		DueDate      time.Time `json:"due_date"`
		Position     int       `json:"position"`
	}

	// Parse request
	var req UpdateCardRequest
	if err := ctx.BodyParser(&req); err != nil {
		return utils.BadRequest(ctx, "Gagal parsing data", err.Error())
	}

	if _, err := uuid.Parse(publicID); err != nil {
		return utils.BadRequest(ctx, "ID tidak valid", err.Error())
	}

	// Tampung data request dalam struct card yang sesuai dengan model card
	card := &models.Card{
		Title:       req.Title,
		Description: req.Description,
		DueDate:     &req.DueDate,
		Position:    req.Position,
		PublicID:    uuid.MustParse(publicID),
	}

	// Update card
	if err := c.service.Update(card, req.ListPublicID); err != nil {
		return utils.InternalServerError(ctx, "Gagal memperbarui card", err.Error())
	}

	return utils.Success(ctx, "Card berhasil diperbarui", card)
}

func (c *CardController) DeleteCard(ctx *fiber.Ctx) error {
	publicID := ctx.Params("id")

	if _, err := uuid.Parse(publicID); err != nil {
		return utils.BadRequest(ctx, "ID tidak valid", err.Error())
	}

	// Ambil data card untuk mengecek apakah card tersebut ada dan mengambil InternalID nya
	card, err := c.service.GetByPublicID(publicID)
	if err != nil {
		return utils.InternalServerError(ctx, "Card tidak ditemukan", err.Error())
	}

	if err := c.service.Delete(uint(card.InternalID)); err != nil {
		return utils.InternalServerError(ctx, "Gagal menghapus card", err.Error())
	}

	return utils.Success(ctx, "Card berhasil dihapus", nil)
}

func (c *CardController) GetCardDetail(ctx *fiber.Ctx) error {
	publicID := ctx.Params("id")

	// Validasi ID
	if _, err := uuid.Parse(publicID); err != nil {
		return utils.BadRequest(ctx, "ID tidak valid", err.Error())
	}

	// Ambil data card
	card, err := c.service.GetByPublicID(publicID)
	if err != nil {
		return utils.InternalServerError(ctx, "Error saat mengambil card", err.Error())
	}

	if card == nil {
		return utils.NotFound(ctx, "Card not found", err.Error())
	}

	return utils.Success(ctx, "Card berhasil diambil", card)
}
