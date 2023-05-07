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

var productCollection *mongo.Collection = configs.GetCollection(configs.DB, "products")

func GetProducts(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	//get all products
	cursor, err := productCollection.Find(ctx, bson.M{})
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.Response{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}

	//iterate over the cursor and add each product to a slice
	products := []models.Product{}
	for cursor.Next(ctx) {
		var product models.Product
		err := cursor.Decode(&product)
		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(responses.Response{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
		}
		products = append(products, product)
	}

	return c.Status(http.StatusOK).JSON(responses.Response{Status: http.StatusOK, Message: "success", Data: &fiber.Map{"products": products}})
}

func GetProductById(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	//get product id from the request params
	id, err := primitive.ObjectIDFromHex(c.Params("id"))
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.Response{Status: http.StatusBadRequest, Message: "error", Data: &fiber.Map{"data": "Bad Request"}})
	}

	//find the product by id
	var product models.Product
	err = productCollection.FindOne(ctx, bson.M{"id": id}).Decode(&product)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.Status(http.StatusNotFound).JSON(responses.Response{Status: http.StatusNotFound, Message: "error", Data: &fiber.Map{"data": "Product not found"}})
		}
		return c.Status(http.StatusInternalServerError).JSON(responses.Response{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}

	return c.Status(http.StatusOK).JSON(responses.Response{Status: http.StatusOK, Message: "success", Data: &fiber.Map{"product": product}})
}

func GetProductsByCategory(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// get category from request params
	category := c.Params("category")

	// get all products with specified category
	cursor, err := productCollection.Find(ctx, bson.M{"category": category})
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.Response{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}

	//iterate over the cursor and add each product to a slice
	products := []models.Product{}
	for cursor.Next(ctx) {
		var product models.Product
		err := cursor.Decode(&product)
		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(responses.Response{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
		}
		products = append(products, product)
	}

	return c.Status(http.StatusOK).JSON(responses.Response{Status: http.StatusOK, Message: "success", Data: &fiber.Map{"products": products}})
}

func AddProduct(c *fiber.Ctx) error {
	// parse request body to product model
	var product models.Product
	err := c.BodyParser(&product)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.Response{Status: http.StatusBadRequest, Message: "error", Data: &fiber.Map{"data": "Bad Request"}})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// insert new product
	result, err := productCollection.InsertOne(ctx, product)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.Response{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}

	// get inserted product by id
	var newProduct models.Product
	err = productCollection.FindOne(ctx, bson.M{"id": result.InsertedID}).Decode(&newProduct)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.Response{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}

	return c.Status(http.StatusOK).JSON(responses.Response{Status: http.StatusOK, Message: "success", Data: &fiber.Map{"product": newProduct}})
}

func EditProduct(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	//get product id from the request params
	id, err := primitive.ObjectIDFromHex(c.Params("id"))
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.Response{Status: http.StatusBadRequest, Message: "error", Data: &fiber.Map{"data": "Bad Request"}})
	}

	// parse request body to product model
	var product models.Product
	err = c.BodyParser(&product)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.Response{Status: http.StatusBadRequest, Message: "error", Data: &fiber.Map{"data": "Bad Request"}})
	}

	// update product by id
	update := bson.M{
		"$set": product,
	}
	_, err = productCollection.UpdateOne(ctx, bson.M{"id": id}, update)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.Response{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}

	// get updated product by id
	var updatedProduct models.Product
	err = productCollection.FindOne(ctx, bson.M{"id": id}).Decode(&updatedProduct)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.Response{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}

	return c.Status(http.StatusOK).JSON(responses.Response{Status: http.StatusOK, Message: "success", Data: &fiber.Map{"product": updatedProduct}})
}

func DeleteProduct(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	//get product id from the request params
	id, err := primitive.ObjectIDFromHex(c.Params("id"))
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.Response{Status: http.StatusBadRequest, Message: "error", Data: &fiber.Map{"data": "Bad Request"}})
	}

	// delete product by id
	_, err = productCollection.DeleteOne(ctx, bson.M{"id": id})
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.Response{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}

	return c.Status(http.StatusOK).JSON(responses.Response{Status: http.StatusOK, Message: "success", Data: &fiber.Map{"message": "Product deleted successfully"}})
}
