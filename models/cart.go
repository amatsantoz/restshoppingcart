package models

import (
	"gorm.io/gorm"
)

type Cart struct {
	gorm.Model
	User_Fk   uint
	Products []Product `gorm:"many2many:cart_products;"`
}

func CreateCart(db *gorm.DB, newCart *Cart, userId uint) (err error) {
	newCart.User_Fk = userId
	err = db.Create(newCart).Error
	if err != nil {
		return err
	}
	return nil
}

func InsertProductToCart(db *gorm.DB, insertedCart *Cart, product *Product) (err error) {
	insertedCart.Products = append(insertedCart.Products, product)
	err = db.Save(insertedCart).Error
	if err != nil {
		return err
	}
	return nil
}

func ReadAllProductsInCart(db *gorm.DB, cart *Cart, id int) (err error) {
	err = db.Where("user_id=?", id).Preload("Products").Find(cart).Error
	if err != nil {
		return err
	}
	return nil
}

func ReadCartById(db *gorm.DB, cart *Cart, id int) (err error) {
	err = db.Where("user_id=?", id).First(cart).Error
	if err != nil {
		return err
	}
	return nil
}