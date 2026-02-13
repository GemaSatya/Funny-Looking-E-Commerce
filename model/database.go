package model

import (
	"fmt"

	"github.com/GemaSatya/E-Commerce/env"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase(){

	connectionUrl := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local",
		env.GetEnv("DB_USERNAME"),
		env.GetEnv("DB_PASSWORD"),
		env.GetEnv("DB_HOST"),
		env.GetEnv("DB_PORT"),
		env.GetEnv("DB_NAME"))

	database, err := gorm.Open(mysql.Open(connectionUrl))
	if err != nil{
		panic("Failed to connect to database!")
	}

	database.AutoMigrate(&User{}, &Product{}, &Order{}, &Cart{}, &CartItem{}, &Category{}, &Login{})

	DB = database

}