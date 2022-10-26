package controllers

import (
	"achmad/restshoppingcart/database"
	"achmad/restshoppingcart/models"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type CartController struct {
	// Declare variables
	Db    *gorm.DB
}

func InitCartController() *CartController {
	db := database.InitDb()
	// gorm sync
	db.AutoMigrate(&models.Cart{})

	return &CartController{Db: db}
}

func (controller *CartController) InsertProductToCart(c *fiber.Ctx) error {
	params := c.AllParams() // "{"id": "1"}"

	intCartId, _ := strconv.Atoi(params["cartid"])
	intProductId, _ := strconv.Atoi(params["productid"])

	var cart models.Cart
	var product models.Product

	// Find the product first, Mencari product dengan idproduct = xx
	if err := models.ReadProductById(controller.Db, &product, intProductId); err != nil {
		// return c.SendStatus(500) // http 500 internal server error
		return c.JSON(fiber.Map{
			"status":  500,
			"message": "Tidak dapat menemukan Product dengan Id " + params["productid"] + ", Gagal menambahkan ke Shopping Cart ",
		})
	}

	// Then find the cart, mencari cart denan idcart = xx
	if err := models.ReadCartById(controller.Db, &cart, intCartId); err != nil {
		// return c.SendStatus(500) // http 500 internal server error
		return c.JSON(fiber.Map{
			"status":  500,
			"message": "Tidak dapat menemukan Cart dengan Id " + params["cartid"] + ", Gagal menambahkan ke Shopping Cart ",
		})
	}

	// Finally, insert the product to cart
	if err := models.InsertProductToCart(controller.Db, &cart, &product); err != nil {
		// return c.SendStatus(500) // http 500 internal server error
		return c.JSON(fiber.Map{
			"status":  500,
			"message": "Internal Server Error, Gagal menambahkan ke Shopping Cart ",
		})
	}

	// if succeed
	return c.JSON(fiber.Map{
		"status":  200,
		"message": "Berhasil Menambahkan Product dengan Id " + params["productid"] + " ke Shopping Cart " + params["cartid"],
	})
}


func (controller *CartController) GetCart(c *fiber.Ctx) error {
	params := c.AllParams() // "{"id": "1"}"

	intCartId, _ := strconv.Atoi(params["cartid"])

	var cart models.Cart
	err := models.ReadAllProductsInCart(controller.Db, &cart, intCartId)
	if err != nil {
		return c.SendStatus(500) // http 500 internal server error
	}

	return c.JSON(fiber.Map{
		"Message":  "Shopping Cart dengan Id " + params["cartid"],
		"Products": cart.Products,
	})
}
