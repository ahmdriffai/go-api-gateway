package main

import (
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/proxy"
	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret = []byte("your-secret-key")

func login(c *fiber.Ctx) error {
	type loginRequest struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	var req loginRequest
	err := c.BodyParser(&req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid request")
	}

	if req.Username == "admin" && req.Password == "admin" {
		token := jwt.NewWithClaims(jwt.SigningMethodES256, jwt.MapClaims{
			"username": req.Username,
			"exp":      time.Now().Add(time.Hour * 2).Unix(),
		})

		tokenString, err := token.SignedString(jwtSecret)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString("Could not generate token")
		}

		return c.JSON(fiber.Map{
			"token": tokenString,
		})
	}

	return c.Status(fiber.StatusUnauthorized).SendString("Invalid credentials")
}

func authMiddleware(c *fiber.Ctx) error {
	tokenString := c.Get("Authorization")

	if tokenString == "" {
		return c.Status(fiber.StatusUnauthorized).SendString("Missing or invalid token")
	}

	// Remove "Bearer " prefix if present
	if len(tokenString) > 7 && tokenString[:7] == "Bearer " {
		tokenString = tokenString[7:]
	}

	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fiber.NewError(fiber.StatusUnauthorized, "Invalid signing method")
		}

		return jwtSecret, nil
	})

	if err != nil || !token.Valid {
		return c.Status(fiber.StatusUnauthorized).SendString("Unauthorized")
	}

	return c.Next()
}

func main() {
	app := fiber.New()

	app.Post("/login", login)

	app.Use(authMiddleware)

	app.All("/tasks/*", func(c *fiber.Ctx) error {
		targetURL := "http://localhost:3000/tasks" + c.OriginalURL()[len("/tasks"):]
		log.Printf("Proxying to Task Service : %s", targetURL)
		return proxy.Do(c, targetURL)
	})

	// Catch-all for unmatched routes
	app.All("*", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusNotFound).SendString("Route not found")
	})

	log.Println("Api gateway runnin on port 8080")
	log.Fatal(app.Listen(":8080"))
}
