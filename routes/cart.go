package routes

import (
	"shoppy/controllers"

	"github.com/gofiber/fiber/v2"
)

func CartRoute(app *fiber.App) {
	// Protected routes
	protected := app.Group("/cart")

	protected.Post("", controllers.AddProduct)
	protected.Get("", controllers.GetCartItems)
	protected.Delete("/:product_id", controllers.DeleteFromCart)
}
