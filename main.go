package main

import (
	"UploadDocumentsAPI/database"
	"UploadDocumentsAPI/middleware"
	"UploadDocumentsAPI/routes"
	"log"
	"os"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/gofiber/fiber/v3/middleware/recover"
	"github.com/joho/godotenv"
)

func main() {

	//Environment variables
	//Se carga el archivo .env con la variable de entorno
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	//Se lee la variable de entorno para saber a que entorno se ejecutará la App
	env := os.Getenv("APP_ENV")
	if "" == env {
		env = "development"
	}

	//Se carga el erchivo con las variables de entorno final
	godotenv.Load(".env." + env)
	log.Print(env)
	port_app := os.Getenv("PORT_APP")
	DB_SERVER := os.Getenv("DB_SERVER")
	DB_NAME := os.Getenv("DB_NAME")
	DB_USER := os.Getenv("DB_USER")
	DB_PASS := os.Getenv("DB_PASS")
	FIREBASE_CREDENTIALS := os.Getenv("FIREBASE_CREDENTIALS")

	//Se inicializa firebase para la verificacion de tokens
	err = middleware.InitFirebase(FIREBASE_CREDENTIALS)
	if err != nil {
		log.Fatal(err)
	}

	//Se inicializa la app con fiber
	app := fiber.New(fiber.Config{
		BodyLimit: 20 * 1024 * 1024, //Set bodyLimit to 10 Mb
	})

	//Se define el origen permitido según el entorno
	allowedOrigins := []string{"http://origos.no-ip.com", "https://origos.no-ip.com"}
	if env == "development" {
		allowedOrigins = []string{"http://localhost:3000", "http://127.0.0.1:3000"}
	}

	//Se inicializa cors de fiber para habilitar Cross-Origin Resource Sharing
	app.Use(cors.New(cors.Config{
		Next:             nil,
		AllowOrigins:     allowedOrigins,
		AllowMethods:     []string{"GET", "POST", "HEAD", "PUT", "DELETE", "PATCH", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Accept", "Content-Type", "Accept-Encoding", "Accept-Language", "Authorization"},
		AllowCredentials: false,
		ExposeHeaders:    []string{},
		MaxAge:           0,
}))

	//database.ConnectMongoDB()
	//Conexion a la BD
	database.ConnectSqlDB(DB_SERVER, DB_NAME, DB_USER, DB_PASS)

	//Manejo de rutas
	routes.Register(app)
	//Se inicializa recover de fiber para el manejo de errores
	app.Use(recover.New())
	log.Println("Server will start at http://localhost:" + port_app)
	log.Fatal(app.Listen("127.0.0.1:" + port_app))

}
