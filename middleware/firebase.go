package middleware

import (
	"context"
	"log"
	"strings"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
	"github.com/gofiber/fiber/v3"
	"google.golang.org/api/option"
)

var firebaseClient *auth.Client

func InitFirebase(credentialsPath string) error {
	ctx := context.Background()
	opt := option.WithCredentialsFile(credentialsPath)
	app, err := firebase.NewApp(ctx, nil, opt)
	if err != nil {
		return err
	}
	firebaseClient, err = app.Auth(ctx)
	if err != nil {
		return err
	}
	log.Println("Firebase Auth initialized")
	return nil
}

func FirebaseAuth(c fiber.Ctx) error {
	header := c.Get("Authorization")
	if header == "" || !strings.HasPrefix(header, "Bearer ") {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "missing or invalid authorization header"})
	}
	tokenString := strings.TrimPrefix(header, "Bearer ")
	token, err := firebaseClient.VerifyIDToken(context.Background(), tokenString)
	if err != nil {
		log.Printf("Firebase token verification failed: %v", err)
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid token"})
	}
	c.Locals("uid", token.UID)
	return c.Next()
}