package routes

import (
	"UploadDocumentsAPI/controllers/uploader"
	"UploadDocumentsAPI/middleware"
	"github.com/gofiber/fiber/v3"
)

func Register(app *fiber.App) {

	contractText := app.Group("/uploader", middleware.FirebaseAuth)
	contractText.Post("/UploadFile", middleware.TenantContext, uploader.UploadFile)
}
