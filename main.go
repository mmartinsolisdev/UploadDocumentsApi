package main

import (
  "UploadDocumentsAPI/database"
  "UploadDocumentsAPI/routes"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {

	//Se inicializa la app con fiber
	app := fiber.New()
	//Se inicializa cors de fiber para habilitar Cross-Origin Resource Sharing
	app.Use(cors.New(cors.Config{
		Next:             nil,
		AllowOrigins:     "*",
		AllowMethods:     "GET,POST,HEAD,PUT,DELETE,PATCH",
		AllowHeaders:  		"Accept, Content-Type, Content-Length, Accept-Encoding, Authorization, application/json",
		AllowCredentials: false,
		ExposeHeaders:    "",
		MaxAge:           0,
}))

//database.ConnectMongoDB()
//Conexion a la BD
database.ConnectSqlDB()

//Manejo de rutas
routes.Register(app)
//Se inicializa recover de fiber para el manejo de errores
	app.Use(recover.New())
	log.Println("Server will start at http://localhost:3000/")
	log.Fatal(app.Listen(":3000"))
}
