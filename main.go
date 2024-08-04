package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Categories struct {
	Id   int    `json:"id" gorm:"primaryKey autoIncrement"`
	Nama string `json:"nama"`
}

type BaseRespose struct {
	Status  bool        `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func main() {
	InitDatabase()
	//	loadEnv()
	e := echo.New()
	e.GET("/categories", GetUsersController)
	e.POST("/categories", AddUsersController)
	e.GET("/categories/:id", GetUserDetailController)
	e.Start(":8000")
}

var DB *gorm.DB

//func loadEnv() {
//	err := godotenv.Load()
//	if err != nil {
//		panic("Failed load env file")
//	}
//}

func InitDatabase() {
	dsn := fmt.Sprintf("%v:%v@tcp(%v:%v)/%v?charset=utf8mb4&parseTime=True&loc=Local",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"))
	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("gagal init database")
	}
	Migration()
}

func Migration() {
	DB.AutoMigrate(&Categories{})
}

func AddUsersController(c echo.Context) error {
	var user Categories
	c.Bind(&user)

	result := DB.Create(&user)

	if result.Error != nil {
		return c.JSON(http.StatusInternalServerError, BaseRespose{
			Status:  false,
			Message: "Failed add data Categories",
			Data:    nil,
		})
	}

	return c.JSON(http.StatusCreated, BaseRespose{
		Status:  true,
		Message: "Success add data Categories",
		Data:    user,
	})
}

func GetUserDetailController(c echo.Context) error {
	id := c.Param("id")

	var users Categories = Categories{1, id}

	return c.JSON(http.StatusOK, BaseRespose{
		Status:  true,
		Message: "Success get detail Categories",
		Data:    users,
	})
}

func GetUsersController(c echo.Context) error {
	var users []Categories

	result := DB.Find(&users)

	if result.Error != nil {
		return c.JSON(http.StatusInternalServerError, BaseRespose{
			Status:  false,
			Message: "Failed get data Categories",
			Data:    nil,
		})
	}

	return c.JSON(http.StatusOK, BaseRespose{
		Status:  true,
		Message: "Success get data Categories",
		Data:    users,
	})
}
