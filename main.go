package main

//import "fmt"
import (
	"log"

	"github.com/gofiber/fiber/v2"
)

func main() {

	//Se inicializa la app con fiber
	app := fiber.New()

  app.Get("/", func(c *fiber.Ctx) error {
  return c.SendString("Hello, World!")
})

	log.Println("Server will start at http://localhost:3000/")
	log.Fatal(app.Listen(":3000"))
}
