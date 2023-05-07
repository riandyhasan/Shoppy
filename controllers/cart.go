package controllers

import (
	"context"
	"net/http"
	"shoppy/configs"
	"shoppy/models"
	"shoppy/responses"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var cartCollection *mongo.Collection = configs.GetCollection(configs.DB, "cart")

func AddToCart(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// get user ID from JWT bearer token
	userID, err := primitive.ObjectIDFromHex(c.Locals("userID").(string))

	//parse the request body to get the product details
	var product_cart models.Add_Products
	if err := c.BodyParser(&product_cart); err != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.Response{Status: http.StatusBadRequest, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}

	//find the product by id
	var product models.Product
	err = productCollection.FindOne(ctx, bson.M{"id": product_cart.Product_ID}).Decode(&product)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.Status(http.StatusNotFound).JSON(responses.Response{Status: http.StatusNotFound, Message: "error", Data: &fiber.Map{"data": "Product not found"}})
		}
		return c.Status(http.StatusInternalServerError).JSON(responses.Response{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}

	//use the validator library to validate required fields
	if validationErr := validate.Struct(&product); validationErr != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.Response{Status: http.StatusBadRequest, Message: "error", Data: &fiber.Map{"data": validationErr.Error()}})
	}

	//create a new cart item
	newCartItem := models.Cart_Products{
		User_ID: userID,
		Product: product,
		Amount:  product_cart.Amount,
	}

	//add the new cart item to the user's cart
	result, err := cartCollection.InsertOne(ctx, newCartItem)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.Response{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}

	return c.Status(http.StatusOK).JSON(responses.Response{Status: http.StatusOK, Message: "success", Data: &fiber.Map{"data": result}})
}

func GetCartItems(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// get user ID from JWT bearer token
	userID, err := primitive.ObjectIDFromHex(c.Locals("userID").(string))
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.Response{Status: http.StatusBadRequest, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}

	// find the user's cart items
	var cartItems []models.Cart_Products
	cur, err := cartCollection.Find(ctx, bson.M{"user_id": userID})
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.Response{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}
	defer cur.Close(ctx)
	for cur.Next(ctx) {
		var cartItem models.Cart_Products
		err := cur.Decode(&cartItem)
		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(responses.Response{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
		}
		cartItems = append(cartItems, cartItem)
	}
	if err := cur.Err(); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.Response{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}

	return c.Status(http.StatusOK).JSON(responses.Response{Status: http.StatusOK, Message: "success", Data: &fiber.Map{"data": cartItems}})
}

func DeleteFromCart(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// get user ID from JWT bearer token
	userID, err := primitive.ObjectIDFromHex(c.Locals("userID").(string))
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.Response{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}

	// parse product ID from request body
	productID, err := primitive.ObjectIDFromHex(c.Params("product_id"))
	if err := c.BodyParser(&productID); err != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.Response{Status: http.StatusBadRequest, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}

	// delete cart item by user ID and product ID
	_, err = cartCollection.DeleteOne(ctx, bson.M{"user_id": userID, "product.id": productID})
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.Response{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}

	return c.Status(http.StatusOK).JSON(responses.Response{Status: http.StatusOK, Message: "success", Data: &fiber.Map{"data": "Product deleted from cart"}})
}
