package controller

import (
	"backend/internal/helpers"
	"backend/internal/modules/auth/dto"
	authservice "backend/internal/modules/auth/service"
	"backend/internal/modules/auth/validator"
	pesertaservice "backend/internal/modules/peserta/service"
	userservice "backend/internal/modules/user/service"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

type AuthController struct {
	authService    authservice.AuthService
	userService    userservice.UserService
	pesertaService pesertaservice.PesertaService
}

func NewAuthController(authService authservice.AuthService, userService userservice.UserService, pesertaService pesertaservice.PesertaService) *AuthController {
	return &AuthController{
		authService:    authService,
		userService:    userService,
		pesertaService: pesertaService,
	}
}

// Register godoc
// @Summary Register user baru
// @Description Mendaftarkan akun user (admin/guru) baru. Untuk akun peserta/siswa, gunakan proses import peserta, bukan endpoint ini.
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body dto.RegisterRequest true "Data registrasi"
// @Success 201 {object} helpers.Response{data=dto.AuthResponse} "Register successfully"
// @Failure 400 {object} helpers.Response "Validasi gagal atau email/username sudah terdaftar"
// @Router /auth/register [post]
func (c *AuthController) Register(ctx *fiber.Ctx) error {
	var req dto.RegisterRequest

	if err := ctx.BodyParser(&req); err != nil {
		return helpers.ErrorResponse(ctx, fiber.StatusBadRequest, "Invalid request format", nil)
	}

	if err := validator.ValidateRegister(req); err != nil {
		return helpers.ErrorResponse(ctx, fiber.StatusBadRequest, "Validation error", nil)
	}

	resp, err := c.authService.Register(&req)
	if err != nil {
		return helpers.ErrorResponse(ctx, fiber.StatusBadRequest, err.Error(), nil)
	}

	return helpers.SuccessResponse(ctx, fiber.StatusCreated, "Register successfully", resp)
}

// Login godoc
// @Summary Login
// @Description Login untuk user (admin/guru) maupun peserta (siswa) menggunakan username & password yang sama. Sistem otomatis mendeteksi apakah username adalah user internal atau peserta (peserta divalidasi ke API sekolah eksternal).
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body dto.LoginRequest true "Kredensial login"
// @Success 200 {object} helpers.Response{data=dto.AuthResponse} "Login successfully, berisi JWT token"
// @Failure 401 {object} helpers.Response "Username tidak ditemukan atau password salah"
// @Router /auth/login [post]
func (c *AuthController) Login(ctx *fiber.Ctx) error {
	var req dto.LoginRequest

	if err := ctx.BodyParser(&req); err != nil {
		return helpers.ErrorResponse(ctx, fiber.StatusBadRequest, "Invalid request format", nil)
	}

	if err := validator.ValidateLogin(req); err != nil {
		return helpers.ErrorResponse(ctx, fiber.StatusBadRequest, "Validation error", nil)
	}

	resp, err := c.authService.Login(&req)
	if err != nil {
		return helpers.ErrorResponse(ctx, fiber.StatusUnauthorized, err.Error(), nil)
	}

	return helpers.SuccessResponse(ctx, fiber.StatusOK, "Login successfully", resp)
}

// GetCurrentUser godoc
// @Summary Info user yang sedang login
// @Description Mengambil data user/peserta yang sedang login berdasarkan JWT token yang dikirim.
// @Tags Auth
// @Produce json
// @Security BearerAuth
// @Success 200 {object} helpers.Response{data=dto.CurrentUserResponse} "Get current user successfully"
// @Failure 401 {object} helpers.Response "Token tidak valid"
// @Failure 404 {object} helpers.Response "User/peserta tidak ditemukan"
// @Router /auth/me [get]
func (c *AuthController) GetCurrentUser(ctx *fiber.Ctx) error {
	user := ctx.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	userID := claims["user_id"].(string)
	role, _ := claims["role"].(string)

	if role == "peserta" {
		pesertaResp, err := c.pesertaService.GetPesertaByID(userID)
		if err != nil {
			return helpers.ErrorResponse(ctx, fiber.StatusNotFound, "Peserta not found", nil)
		}
		return helpers.SuccessResponse(ctx, fiber.StatusOK, "Get current user successfully", dto.CurrentUserResponse{
			ID:    pesertaResp.ID,
			Name:  pesertaResp.Nama,
			Email: pesertaResp.Username,
			Role:  role,
		})
	}

	userResp, err := c.userService.GetByID(userID)
	if err != nil {
		return helpers.ErrorResponse(ctx, fiber.StatusNotFound, "User not found", nil)
	}

	return helpers.SuccessResponse(ctx, fiber.StatusOK, "Get current user successfully", dto.CurrentUserResponse{
		ID:    userResp.ID,
		Name:  userResp.Name,
		Email: userResp.Email,
		Role:  userResp.Role,
	})
}
