package middleware

import (
	"errors"
	"strings"

	"UploadDocumentsAPI/database"

	"github.com/gofiber/fiber/v3"
)

// TenantContext resuelve la base de datos del cliente a partir de los headers
// X-Client (obligatorio) y X-Database (opcional; vacío ⇒ base por defecto del cliente).
func TenantContext(c fiber.Ctx) error {
	slug := c.Get("X-Client")
	if slug == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "missing X-Client header"})
	}

	databaseName := strings.TrimSpace(c.Get("X-Database"))
	db, err := database.GetDB(slug, databaseName)
	if err != nil {
		switch {
		case errors.Is(err, database.ErrUnknownTenant):
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "unknown client"})
		case errors.Is(err, database.ErrUnknownDatabase):
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "unknown database"})
		default:
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "database connection error"})
		}
	}

	c.Locals("db", db)
	return c.Next()
}
