package middleware

import (
	"strings"

	"backend/configs"
	"backend/internal/helpers"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func JWTAuth() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		authHeader := ctx.Get("Authorization")
		if authHeader == "" {
			return helpers.ErrorResponse(ctx, fiber.StatusUnauthorized, "Missing authorization header", nil)
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			return helpers.ErrorResponse(ctx, fiber.StatusUnauthorized, "Invalid authorization format", nil)
		}

		tokenString := parts[1]
		jwtConfig := configs.GetJWTConfig()

		token, err := jwt.ParseWithClaims(tokenString, jwt.MapClaims{}, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fiber.NewError(fiber.StatusUnauthorized, "Invalid signing method")
			}
			return []byte(jwtConfig.Secret), nil
		})

		if err != nil || !token.Valid {
			return helpers.ErrorResponse(ctx, fiber.StatusUnauthorized, "Invalid token", nil)
		}

		ctx.Locals("user", token)
		return ctx.Next()
	}
}

// JWTOptional mencoba parse JWT tapi tidak memblokir jika tidak ada/invalid.
// user_id akan tersedia di ctx.Locals("user_id") jika token valid.
func JWTOptional() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		authHeader := ctx.Get("Authorization")
		if authHeader == "" {
			return ctx.Next()
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			return ctx.Next()
		}

		jwtConfig := configs.GetJWTConfig()
		token, err := jwt.ParseWithClaims(parts[1], jwt.MapClaims{}, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fiber.NewError(fiber.StatusUnauthorized, "Invalid signing method")
			}
			return []byte(jwtConfig.Secret), nil
		})

		if err == nil && token.Valid {
			if claims, ok := token.Claims.(jwt.MapClaims); ok {
				if userID, ok := claims["user_id"].(string); ok {
					ctx.Locals("user_id", userID)
				}
			}
		}

		return ctx.Next()
	}
}
