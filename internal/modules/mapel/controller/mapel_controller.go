package controller

import (
	"strconv"

	"backend/internal/helpers"
	"backend/internal/modules/mapel/dto"
	"backend/internal/modules/mapel/service"

	"github.com/gofiber/fiber/v2"
)

type MapelController struct {
	service service.MapelService
}

func NewMapelController(service service.MapelService) *MapelController {
	return &MapelController{service: service}
}

// CreateMapel godoc
// @Summary Buat mata pelajaran baru
// @Tags Mapel
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.CreateMapelRequest true "Data mapel"
// @Success 201 {object} helpers.Response{data=dto.MapelResponse} "Create mapel successfully"
// @Failure 400 {object} helpers.Response "Validasi gagal"
// @Router /mapel [post]
func (c *MapelController) CreateMapel(ctx *fiber.Ctx) error {
	var req dto.CreateMapelRequest

	if err := ctx.BodyParser(&req); err != nil {
		return helpers.ErrorResponse(ctx, fiber.StatusBadRequest, "Invalid request format", nil)
	}

	resp, err := c.service.CreateMapel(&req)
	if err != nil {
		return helpers.ErrorResponse(ctx, fiber.StatusBadRequest, err.Error(), nil)
	}

	return helpers.SuccessResponse(ctx, fiber.StatusCreated, "Create mapel successfully", resp)
}

// GetAllMapel godoc
// @Summary List mata pelajaran
// @Tags Mapel
// @Produce json
// @Param page query int false "Halaman" default(1)
// @Param page_size query int false "Jumlah per halaman" default(10)
// @Success 200 {object} helpers.Response{data=dto.MapelListResponse} "Get all mapel successfully"
// @Failure 500 {object} helpers.Response "Gagal mengambil data"
// @Router /mapel [get]
func (c *MapelController) GetAllMapel(ctx *fiber.Ctx) error {
	page := ctx.Query("page", "1")
	pageSize := ctx.Query("page_size", "10")

	pageNum, err := strconv.Atoi(page)
	if err != nil || pageNum <= 0 {
		pageNum = 1
	}

	pageSizeNum, err := strconv.Atoi(pageSize)
	if err != nil || pageSizeNum <= 0 {
		pageSizeNum = 10
	}

	resp, err := c.service.GetAllMapel(pageNum, pageSizeNum)
	if err != nil {
		return helpers.ErrorResponse(ctx, fiber.StatusInternalServerError, err.Error(), nil)
	}

	return helpers.SuccessResponse(ctx, fiber.StatusOK, "Get all mapel successfully", resp)
}

// GetMapelByID godoc
// @Summary Detail mata pelajaran
// @Tags Mapel
// @Produce json
// @Param id path string true "ID mapel (uuid)"
// @Success 200 {object} helpers.Response{data=dto.MapelResponse} "Get mapel successfully"
// @Failure 404 {object} helpers.Response "Mapel tidak ditemukan"
// @Router /mapel/{id} [get]
func (c *MapelController) GetMapelByID(ctx *fiber.Ctx) error {
	id := ctx.Params("id")

	resp, err := c.service.GetMapelByID(id)
	if err != nil {
		return helpers.ErrorResponse(ctx, fiber.StatusNotFound, err.Error(), nil)
	}

	return helpers.SuccessResponse(ctx, fiber.StatusOK, "Get mapel successfully", resp)
}

// UpdateMapel godoc
// @Summary Update mata pelajaran
// @Tags Mapel
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "ID mapel (uuid)"
// @Param request body dto.UpdateMapelRequest true "Data mapel"
// @Success 200 {object} helpers.Response{data=dto.MapelResponse} "Update mapel successfully"
// @Failure 400 {object} helpers.Response "Validasi gagal"
// @Router /mapel/{id} [put]
func (c *MapelController) UpdateMapel(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	var req dto.UpdateMapelRequest

	if err := ctx.BodyParser(&req); err != nil {
		return helpers.ErrorResponse(ctx, fiber.StatusBadRequest, "Invalid request format", nil)
	}

	resp, err := c.service.UpdateMapel(id, &req)
	if err != nil {
		return helpers.ErrorResponse(ctx, fiber.StatusBadRequest, err.Error(), nil)
	}

	return helpers.SuccessResponse(ctx, fiber.StatusOK, "Update mapel successfully", resp)
}

// DeleteMapel godoc
// @Summary Hapus mata pelajaran (soft delete)
// @Tags Mapel
// @Produce json
// @Security BearerAuth
// @Param id path string true "ID mapel (uuid)"
// @Success 200 {object} helpers.Response "Delete mapel successfully"
// @Failure 400 {object} helpers.Response "Gagal menghapus"
// @Router /mapel/{id} [delete]
func (c *MapelController) DeleteMapel(ctx *fiber.Ctx) error {
	id := ctx.Params("id")

	err := c.service.DeleteMapel(id)
	if err != nil {
		return helpers.ErrorResponse(ctx, fiber.StatusBadRequest, err.Error(), nil)
	}

	return helpers.SuccessResponse(ctx, fiber.StatusOK, "Delete mapel successfully", nil)
}

// RestoreMapel godoc
// @Summary Restore mata pelajaran yang sudah dihapus
// @Tags Mapel
// @Produce json
// @Security BearerAuth
// @Param id path string true "ID mapel (uuid)"
// @Success 200 {object} helpers.Response "Restore mapel successfully"
// @Failure 400 {object} helpers.Response "Gagal restore"
// @Router /mapel/{id}/restore [patch]
func (c *MapelController) RestoreMapel(ctx *fiber.Ctx) error {
	id := ctx.Params("id")

	err := c.service.RestoreMapel(id)
	if err != nil {
		return helpers.ErrorResponse(ctx, fiber.StatusBadRequest, err.Error(), nil)
	}

	return helpers.SuccessResponse(ctx, fiber.StatusOK, "Restore mapel successfully", nil)
}
