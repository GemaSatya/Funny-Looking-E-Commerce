package model

import "gorm.io/gorm"

type Product struct {
	gorm.Model
	Name        string
	Description string
	Price       float64
}

type User struct {
	gorm.Model
	Username string
	Email    string
	Password string
}

type Cart struct {
	gorm.Model
	UserID uint
	User   User       `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
	Items  []CartItem `gorm:"foreignKey:CartID;constraint:OnDelete:CASCADE"`
}

type CartItem struct {
	gorm.Model
	CartID    uint
	ProductID uint
	Quantity  int
	Cart      Cart    `gorm:"foreignKey:CartID"`
	Product   Product `gorm:"foreignKey:ProductID"`
}

type Order struct {
	gorm.Model
	UserID    uint
	ProductID uint
	Quantity  int
}

type Category struct {
	gorm.Model
	Name string
}