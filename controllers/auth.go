package controllers

import (
	"achmad/restshoppingcart/database"
	"achmad/restshoppingcart/models"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type LoginForm struct {
	Username string `form:"username" json:"username" validate:"required"`
	Password string `form:"password" json:"password" validate:"required"`
}

type AuthController struct {
	// Declare variables
	Db    *gorm.DB
}

func InitAuthController() *AuthController {
	db := database.InitDb()
	// gorm sync
	db.AutoMigrate(&models.User{})

	return &AuthController{Db: db}
}


// post /login
func (controller *AuthController) LoginPosted(c *fiber.Ctx) error {

	var user models.User
	var myform LoginForm

	if err := c.BodyParser(&myform); err != nil {
		return c.SendStatus(400) // Bad Request, RegisterForm is not complete
	}
	
	errs := models.FindUserByUsername(controller.Db, &user, myform.Username)
	if errs != nil {
		return c.JSON(fiber.Map{
		"message": "Username tidak sama",
	}) // Unsuccessful login (cannot find user)
	}

	compare := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(myform.Password))

	// Throws Unauthorized error
	if myform.Username != user.Username || compare != nil {
		return c.SendStatus(fiber.StatusUnauthorized)
	}
	
// Create the Claims
	exp := time.Now().Add(time.Hour * 72)
	claims := jwt.MapClaims{
		"username": myform.Username,
		"admin": true,
		"exp":   exp.Unix(),
	}

	// Create token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Generate encoded token and send it as response.
	t, err := token.SignedString([]byte("mysecretpassword"))
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	return c.JSON(fiber.Map{
		"message": "Berhasil Login",
		"token": t,
		"expired": exp.Format("2006-01-02 15:04:05"),
	})
}



// POST /register
func (controller *AuthController) AddRegisteredUser(c *fiber.Ctx) error {
	var user models.User
	var cart models.Cart

	if err := c.BodyParser(&user); err != nil {
		return c.SendStatus(400) // Bad Request, RegisterForm is not complete
	}

	// Hash password
	bytes, _ := bcrypt.GenerateFromPassword([]byte(user.Password), 10)
	sHash := string(bytes)

	// Simpan hashing, bukan plain passwordnya
	user.Password = sHash

	// save user
	err := models.CreateUser(controller.Db, &user)
	if err != nil {
		return c.SendStatus(500) // Server error, gagal menyimpan user
	}

	// also create cart
	errCart := models.CreateCart(controller.Db, &cart, user.ID)
	if errCart != nil {
		// Server error, gagal menyimpan user
		return c.JSON(fiber.Map{
			"status":  500,
			"message": "Server error, gagal menyimpan user",
		})
	}

	// if succeed
	return c.JSON(user)
}