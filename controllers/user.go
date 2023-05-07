package controllers

import (
	"context"
	"net/http"
	"shoppy/configs"
	"shoppy/models"
	"shoppy/responses"
	"shoppy/utils"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var userCollection *mongo.Collection = configs.GetCollection(configs.DB, "users")
var validate = validator.New()

func Register(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var user models.User
	defer cancel()

	//validate the request body
	if err := c.BodyParser(&user); err != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.Response{Status: http.StatusBadRequest, Message: "error", Data: &fiber.Map{"data": "Bad Request"}})
	}

	//use the validator library to validate required fields
	if validationErr := validate.Struct(&user); validationErr != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.Response{Status: http.StatusBadRequest, Message: "error", Data: &fiber.Map{"data": "Bad Request"}})
	}

	//check if the user exists in the database
	var db_user models.User
	err := userCollection.FindOne(ctx, bson.M{"email": user.Email}).Decode(&db_user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			// hash password
			hashedPassword, err := utils.HashPassword(user.Password)
			if err != nil {
				return c.Status(http.StatusInternalServerError).JSON(responses.Response{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
			}

			newUser := models.User{
				Id:       primitive.NewObjectID(),
				Name:     user.Name,
				Email:    user.Email,
				Password: hashedPassword,
			}
			result, err := userCollection.InsertOne(ctx, newUser)
			if err != nil {
				return c.Status(http.StatusInternalServerError).JSON(responses.Response{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})

			}
			return c.Status(http.StatusCreated).JSON(responses.Response{Status: http.StatusCreated, Message: "success", Data: &fiber.Map{"data": result}})
		}
		return c.Status(http.StatusInternalServerError).JSON(responses.Response{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}

	return c.Status(http.StatusUnauthorized).JSON(responses.Response{Status: http.StatusUnauthorized, Message: "error", Data: &fiber.Map{"data": "Email has been taken"}})
}

func Login(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	//parse the request body to get user credentials
	var creds models.Authentication
	if err := c.BodyParser(&creds); err != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.Response{Status: http.StatusBadRequest, Message: "error", Data: &fiber.Map{"data": "Bad Request"}})
	}

	//use the validator library to validate required fields
	if validationErr := validate.Struct(&creds); validationErr != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.Response{Status: http.StatusBadRequest, Message: "error", Data: &fiber.Map{"data": "Bad Request"}})
	}

	//check if the user exists in the database
	var user models.User
	err := userCollection.FindOne(ctx, bson.M{"email": creds.Email}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.Status(http.StatusUnauthorized).JSON(responses.Response{Status: http.StatusUnauthorized, Message: "error", Data: &fiber.Map{"data": "Invalid email or password"}})
		}
		return c.Status(http.StatusInternalServerError).JSON(responses.Response{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}

	//check if the password matches
	if !utils.CheckPasswordHash(creds.Password, user.Password) {
		return c.Status(http.StatusUnauthorized).JSON(responses.Response{Status: http.StatusUnauthorized, Message: "error", Data: &fiber.Map{"data": "Invalid email or password"}})
	}

	//create JWT token
	claims := jwt.MapClaims{
		"user_id": user.Id,
		"exp":     time.Now().Add(time.Hour * 24).Unix(), // expiration time
	}
	token, err := utils.CreateToken(claims)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.Response{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}

	return c.Status(http.StatusOK).JSON(responses.Response{Status: http.StatusOK, Message: "success", Data: &fiber.Map{"token": token}})
}
