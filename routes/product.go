package routes

import (
	"shoppy/controllers"

	"github.com/gofiber/fiber/v2"
)

func ProductRoute(app *fiber.App) {
	app.Get("/product", controllers.GetProducts)
	app.Get("/product/:id", controllers.GetProductById)
	app.Get("/product/category/:category", controllers.GetProductsByCategory)
	app.Post("/product", controllers.AddProduct)
	app.Put("/producs/:id", controllers.EditProduct)
}
