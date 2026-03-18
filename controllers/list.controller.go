package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/kasyaproject/sistem-project-management/models"
	"github.com/kasyaproject/sistem-project-management/services"
	"github.com/kasyaproject/sistem-project-management/utils"
)

type ListController struct {
	service services.ListService
}

func NewListController(s services.ListService) *ListController {
	return &ListController{service: s}
}

func (c *ListController) CreateList(ctx *fiber.Ctx) error {
	// Deklarasi variable
	list := new(models.List)

	// Check input data
	if err := ctx.BodyParser(list); err != nil {
		return utils.BadRequest(ctx, "Gagal parsing data", err.Error())
	}

	// Create list
	if err := c.service.Create(list); err != nil {
		return utils.BadRequest(ctx, "Gagal membuat list", err.Error())
	}

	return utils.Success(ctx, "List berhasil dibuat", list)
}

func (c *ListController) UpdateList(ctx *fiber.Ctx) error {
	publicID := ctx.Params("id")
	list := new(models.List)

	// Check input data
	if err := ctx.BodyParser(list); err != nil {
		return utils.BadRequest(ctx, "Gagal parsing data", err.Error())
	}

	// Validasi ID
	if _, err := uuid.Parse(publicID); err != nil {
		return utils.BadRequest(ctx, "ID tidak valid", err.Error())
	}

	// Ambil data list
	existingList, err := c.service.GetByPublicID(publicID)
	if err != nil {
		return utils.NotFound(ctx, "List not found", err.Error())
	}

	list.InternalID = existingList.InternalID // set internal id sesuai dengan list yang sudah ada
	list.PublicID = existingList.PublicID     // set public id sesuai dengan list yang sudah ada

	// Update list
	if err := c.service.Update(list); err != nil {
		return utils.BadRequest(ctx, "Gagal update list", err.Error())
	}

	// Ambil data list yang sudah di update
	listUpdated, err := c.service.GetByPublicID(publicID)
	if err != nil {
		return utils.NotFound(ctx, "Gagal ambil data list yang sudah di update", err.Error())
	}

	return utils.Success(ctx, "List berhasil diupdate", listUpdated)
}

func (c *ListController) GetListOnBoard(ctx *fiber.Ctx) error {
	boardPublicID := ctx.Params("board_id")

	// Validasi ID
	if _, err := uuid.Parse(boardPublicID); err != nil {
		return utils.BadRequest(ctx, "ID tidak valid", err.Error())
	}

	lists, err := c.service.GetByBoardID(boardPublicID)
	if err != nil {
		return utils.NotFound(ctx, "List not found", err.Error())
	}

	return utils.Success(ctx, "List berhasil diambil", lists)
}

func (c *ListController) DeleteList(ctx *fiber.Ctx) error {
	publicID := ctx.Params("id")

	// Validasi ID
	if _, err := uuid.Parse(publicID); err != nil {
		return utils.BadRequest(ctx, "ID tidak valid", err.Error())
	}

	// Cek apakah list ada di board
	list, err := c.service.GetByPublicID(publicID)
	if err != nil {
		return utils.NotFound(ctx, "List not found", err.Error())
	}

	// Delete list
	if err := c.service.Delete(uint(list.InternalID)); err != nil {
		return utils.BadRequest(ctx, "Gagal delete list", err.Error())
	}

	return utils.Success(ctx, "List berhasil dihapus", nil)
}
