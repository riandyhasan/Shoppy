package middlewares

import (
	"fmt"
	"net/http"
	"shoppy/configs"

	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
)

var JWT_SECRET = configs.EnvJWTSecret()

func AuthMiddleware(next fiber.Handler) fiber.Handler {
	return func(c *fiber.Ctx) error {
		tokenString := c.Get("Authorization")
		if tokenString == "" {
			return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"message": "Unauthorized"})
		}

		// Validate the token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Check the signing method
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("invalid token signing method")
			}

			// Get the secret key from the environment variable
			secret := JWT_SECRET

			return []byte(secret), nil
		})
		if err != nil {
			return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"message": "Unauthorized"})
		}

		// Check if the token is valid
		if !token.Valid {
			return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"message": "Unauthorized"})
		}

		// Pass the request to the next handler
		return next(c)
	}
}

// func Protected() fiber.HandlerFunc {
// 	return func(c *fiber.Ctx) error {
// 		// Call the AuthMiddleware to check the authentication token in the Authorization header
// 		return AuthMiddleware(func(c *fiber.Ctx) error {
// 			return c.Next()
// 		})(c)
// 	}
// }
