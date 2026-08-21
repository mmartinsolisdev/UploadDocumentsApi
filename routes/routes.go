package routes

import (

	"UploadDocumentsAPI/controllers/uploader"
	"UploadDocumentsAPI/middleware"
	"github.com/gofiber/fiber/v2"
)

func Register(app *fiber.App) {

	contractText := app.Group("/uploader", middleware.FirebaseAuth)
	contractText.Get("/MembershipsList", uploader.GetMembershipsList)
	contractText.Get("/CombosList", uploader.GetCombos)
	contractText.Post("/UploadFile", uploader.UploadFile)
}
