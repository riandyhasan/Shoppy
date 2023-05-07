package routes

import (
	"shoppy/controllers"

	"github.com/gofiber/fiber/v2"
)

func UserRoute(app *fiber.App) {
	app.Post("/register", controllers.Register)
	app.Post("/login", controllers.Login)
}
