package controllers

import (
	"achmad/restshoppingcart/database"
	"achmad/restshoppingcart/models"
	"fmt"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type TransaksiController struct {
	// Declare variables
	Db    *gorm.DB
}

func InitTransaksiController() *TransaksiController {
	db := database.InitDb()
	// gorm sync
	db.AutoMigrate(&models.Transaksi{})

	return &TransaksiController{Db: db}
}



func (controller *TransaksiController) InsertToTransaksi(c *fiber.Ctx) error {
	params := c.AllParams() // "{"id": "1"}"

	intUserId, _ := strconv.Atoi(params["userid"])

	var transaksi models.Transaksi
	var cart models.Cart

	// Find the cart
	if err := models.ReadCartById(controller.Db, &cart, intUserId); err != nil {
		return c.SendStatus(500) // http 500 internal server error
	}

	//searcp products in cart
	if err := models.ReadAllProductsInCart(controller.Db, &cart, intUserId); err != nil {
		return c.SendStatus(500) // http 500 internal server error
	}

	fmt.Println(cart.Products)
	// Jika Cart kosong
	if len(cart.Products) == 0 {
		return c.JSON(fiber.Map{
			"status":  400,
			"message": "Cart kosong, silahkan isi Product ke Cart terlebih dahulu",
		})
	}

	//membuat transaksi
	if err := models.CreateTransaksi(controller.Db, &transaksi, uint(intUserId), cart.Products); err != nil {
		return c.SendStatus(500) // http 500 internal server error
	}

	// Delete products in cart
	if err := models.UpdateCart(controller.Db, cart.Products, &cart, uint(intUserId));err != nil {
		return c.SendStatus(500) // http 500 internal server error
	}

	// if succeed
	return c.JSON(fiber.Map{
		"status":  200,
		"message": "Berhasil Melakukan Checkout",
	})
}


func (controller *TransaksiController) GetTransaksi(c *fiber.Ctx) error {
	params := c.AllParams() // "{"id": "1"}"

	intUserId, _ := strconv.Atoi(params["userid"])

	var transaksis []models.Transaksi
	err := models.ReadTransaksiById(controller.Db, &transaksis, intUserId)
	if err != nil {
		return c.SendStatus(500) // http 500 internal server error
	}
	return c.JSON(fiber.Map{
		"Title":      "History Transaksi",
		"Transaksis": transaksis,
	})

}


func (controller *TransaksiController) DetailTransaksi(c *fiber.Ctx) error {
	params := c.AllParams() // "{"id": "1"}"
	intTransaksiId, _ := strconv.Atoi(params["transaksiid"])
	// param := struct {ID uint `params:"transaksiid"`}{}
  // c.ParamsParser(&param)

	var transaksi models.Transaksi
	err := models.ReadAllProductsInTransaksi(controller.Db, &transaksi, intTransaksiId)
	if err != nil {
		return c.SendStatus(500) // http 500 internal server error
	}
	return c.JSON(fiber.Map{
		"Title":    "History Transaksi",
		"Products": transaksi.Products,
	})
}
