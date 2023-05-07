package main

import (
	"shoppy/configs"
	"shoppy/routes"

	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()

	// DB
	configs.ConnectDB()

	// Routes
	routes.UserRoute(app)

	// Protected routes
	routes.ProductRoute(app)
	routes.CartRoute(app)
	// routes.TransactionRoute(app)

	app.Listen(":8080")
}
