package main

import (
	"achmad/restshoppingcart/controllers"

	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v3"
)

func main() {
	// session
	// store := session.New()

	// load template engine
	// engine := html.New("./views", ".html")

	app := fiber.New()

	// static
	app.Static("/", "./public", fiber.Static{
		Index: "",
	})
	
	authController := controllers.InitAuthController()
	prodController := controllers.InitProductController()
	cartController := controllers.InitCartController()
	transaksiController := controllers.InitTransaksiController()
	
	app.Post("/register", authController.AddRegisteredUser)
	app.Post("/login", authController.LoginPosted)
	
	prod := app.Group("/products")
	prod.Get("/", prodController.GetAllProduct)
	prod.Get("/detail/:id", prodController.DetailProduct)
	
	app.Use(jwtware.New(jwtware.Config{
		SigningKey: []byte("mysecretpassword"),
	}))
	
	prod.Post("/create",  prodController.AddPostedProduct)
	prod.Post("/ubah/:id",  prodController.AddUpdatedProduct)
	prod.Get("/hapus/:id",  prodController.DeleteProduct)
	prod.Get("/addtocart/:cartid/product/:productid",  cartController.InsertProductToCart)
	
	cart := app.Group("/shoppingcart")
	cart.Get("/:cartid",  cartController.GetCart)
	
	checkout := app.Group("/checkout")
	checkout.Get("/:userid", transaksiController.InsertToTransaksi)
	
	transaksi := app.Group("/transaksi")
	transaksi.Get("/:userid",  transaksiController.GetTransaksi)
	transaksi.Get("/detail/:transaksiid",  transaksiController.DetailTransaksi)
	
	// Middleware to check login
	// CheckLogin := func(c *fiber.Ctx) error {
	// 	sess, _ := store.Get(c)
	// 	val := sess.Get("username")
	// 	if val != nil {
	// 		return c.Next()
	// 	}

	// 	return c.Redirect("/login")
	// }

	// controllers

	// prod.Get("/create",  prodController.AddProduct)
	// prod.Get("/ubah/:id", prodController.UpdateProduct)




	// app.Get("/login", authController.Login)
	// app.Get("/logout", authController.Logout)
	// app.Get("/register", authController.Register)

	app.Listen(":3000")
}
