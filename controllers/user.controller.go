package controllers

import (
	"math"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/jinzhu/copier"
	"github.com/kasyaproject/sistem-project-management/models"
	"github.com/kasyaproject/sistem-project-management/services"
	"github.com/kasyaproject/sistem-project-management/utils"
)

type UserController struct {
	service services.UserService
}

func NewUserController(s services.UserService) *UserController {
	return &UserController{service: s}
}

// Register Controller
func (c *UserController) Register(ctx *fiber.Ctx) error {
	// Struct input data
	user := new(models.User)

	// Check input data
	if err := ctx.BodyParser(user); err != nil {
		return utils.BadRequest(ctx, "Gagal Parsing data", err.Error())
	}

	// Check user register
	if err := c.service.Register(user); err != nil {
		return utils.BadRequest(ctx, "Gagal Registrasi", err.Error())
	}

	// Take out Credential data user (password, internalId) sesuai dengan struct UserResponse di models user
	var userResponse models.UserResponse
	_ = copier.Copy(&userResponse, &user)

	return utils.Success(ctx, "Register Success", userResponse)
}

// Login Controller
func (c *UserController) Login(ctx *fiber.Ctx) error {
	// Struct input data
	var body struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	// Check input data
	if err := ctx.BodyParser(&body); err != nil {
		return utils.BadRequest(ctx, "Invalid Request", err.Error())
	}

	// Check user login
	user, err := c.service.Login(body.Email, body.Password)
	if err != nil {
		return utils.Unauthorized(ctx, "Login Failed!", err.Error())
	}

	// Generate Credential
	token, _ := utils.GenerateToken(user.InternalID, user.Role, user.Email, user.PublicID)
	refreshToken, _ := utils.GenerateRefreshToken(user.InternalID)

	// Take out Credential data user (password, internalId) sesuai dengan struct UserResponse di models user
	var userResponse models.UserResponse
	_ = copier.Copy(&userResponse, &user)

	return utils.Success(ctx, "Login Success", fiber.Map{
		"access_token":  token,
		"refresh_token": refreshToken,
		"user":          userResponse,
	})
}

// Get myprofile Controller
func (c *UserController) GetUser(ctx *fiber.Ctx) error {
	id := ctx.Params("id")

	user, err := c.service.GetByPublicID(id)

	if err != nil {
		return utils.NotFound(ctx, "Data not found", err.Error())
	}

	// Take out Credential data user (password, internalId) sesuai dengan struct UserResponse di models user
	var userResponse models.UserResponse
	err = copier.Copy(&userResponse, &user)

	if err != nil {
		return utils.BadRequest(ctx, "Internal server Error", err.Error())
	}

	return utils.Success(ctx, "Get Data Successfull", userResponse)
}

// Find All user with pagination and params
func (c *UserController) FindAllUser(ctx *fiber.Ctx) error {
	page, _ := strconv.Atoi(ctx.Query("page", "1"))
	limit, _ := strconv.Atoi(ctx.Query("limit", "10"))
	offset := (page - 1) * limit

	filter := ctx.Query("filter", "")
	sort := ctx.Query("sort", "")

	users, total, err := c.service.FindAllUser(filter, sort, limit, offset)
	if err != nil {
		return utils.BadRequest(ctx, "Gagal mengambil data", err.Error())
	}

	// Take out Credential data user (password, internalId) sesuai dengan struct UserResponse di models user
	var userResponse []models.UserResponse
	_ = copier.Copy(&userResponse, &users)

	// buat output data untuk pagination kedalam response
	meta := utils.PaginationMeta{
		Page:      page,
		Limit:     limit,
		Total:     int(total),
		TotalPage: int(math.Ceil(float64(total) / float64(limit))),

		Filter: filter,
		Sort:   sort,
	}

	if total == 0 {
		return utils.NotFoundPagination(ctx, "Data pengguna tidak ditemukan", userResponse, meta)
	}

	return utils.SuccessPagination(ctx, "Data Berhasil Diambil", userResponse, meta)
}

func (c *UserController) UpdateUser(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	publicID, err := uuid.Parse(id)
	if err != nil {
		return utils.BadRequest(ctx, "Invalid ID Format", err.Error())
	}

	// Struct data user dan menyeseuaikan body dengan struct
	var user models.User
	if err := ctx.BodyParser(&user); err != nil {
		return utils.BadRequest(ctx, "Gagal Parsing data", err.Error())
	}

	// Update user
	user.PublicID = publicID
	if err := c.service.Update(&user); err != nil {
		return utils.BadRequest(ctx, "Gagal Update data", err.Error())
	}

	// Ambil data user yang sudah di update
	userUpdated, err := c.service.GetByPublicID(id)
	if err != nil {
		return utils.InternalServerError(ctx, "Gagal Ambil data", err.Error())
	}

	// Take out Credential data user (password, internalId) sesuai dengan struct UserResponse di models user
	var userResponse models.UserResponse
	err = copier.Copy(&userResponse, &userUpdated)
	if err != nil {
		return utils.InternalServerError(ctx, "Gagal Parsing data response", err.Error())
	}

	// return response
	return utils.Success(ctx, "Data Berhasil Diupdate", userResponse)
}

func (c *UserController) DeleteUser(ctx *fiber.Ctx) error {
	id, _ := strconv.Atoi(ctx.Params("id"))
	if err := c.service.Delete(uint(id)); err != nil {
		return utils.InternalServerError(ctx, "Gagal Delete data", err.Error())
	}

	return utils.Success(ctx, "Data Berhasil Dihapus", nil)
}
