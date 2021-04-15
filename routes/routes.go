package routes

import (

	"UploadDocumentsAPI/controllers/uploader"
	"github.com/gofiber/fiber/v2"
)

func Register(app *fiber.App) {

	contractText := app.Group("/uploader", func(c *fiber.Ctx) error {
		return c.Next()
	})
	contractText.Get("/MembershipsList", uploader.GetMembershipsList)
	contractText.Get("/CombosList", uploader.GetCombos)
	contractText.Post("/UploadFile", uploader.UploadFile)
}
