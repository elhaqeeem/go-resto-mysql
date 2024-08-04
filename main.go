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
	Id      int       `json:"id" gorm:"primaryKey autoIncrement"`
	Nama    string    `json:"nama"`
	Minuman []Minuman `json:"minuman" gorm:"foreignKey:Categories_id"`
	Makanan []Makanan `json:"makanan" gorm:"foreignKey:Categories_id"`
}

type Minuman struct {
	ID            uint       `json:"id" gorm:"primaryKey"`
	Nama          string     `json:"nama" gorm:"size:50;not null"`
	Varian        string     `json:"varian" gorm:"size:50"`
	Harga         float64    `json:"harga" gorm:"not null;type:decimal(10,2)"`
	Categories_id uint       `json:"categories_id"` // Foreign Key
	Categories    Categories `json:"categories" gorm:"foreignKey:Categories_id;constraint:OnDelete:CASCADE;"`
}

type Makanan struct {
	ID            uint       `json:"id" gorm:"primaryKey"`
	Nama          string     `json:"nama" gorm:"size:50;not null"`
	Varian        string     `json:"varian" gorm:"size:50"`
	Harga         float64    `json:"harga" gorm:"not null;type:decimal(10,2)"`
	Categories_id uint       `json:"categories_id"` // Foreign Key
	Categories    Categories `json:"categories" gorm:"foreignKey:Categories_id;constraint:OnDelete:CASCADE;"`
}

type Promo struct {
	ID    uint    `json:"id" gorm:"primaryKey"`
	Nama  string  `json:"nama" gorm:"size:100;not null"`
	Harga float64 `json:"harga" gorm:"not null;type:decimal(10,2)"`
}

type Printer struct {
	ID   string `json:"id" gorm:"primaryKey;size:1"`
	Nama string `json:"nama" gorm:"size:50;not null"`
}

type Meja struct {
	ID    uint `json:"id" gorm:"primaryKey"`
	Nomor int  `json:"nomor" gorm:"unique;not null"`
}

type Orders struct {
	ID     uint `json:"id" gorm:"primaryKey"`
	MejaID uint `json:"meja_id" gorm:"not null"`
	Meja   Meja `json:"meja" gorm:"foreignKey:MejaID"`
}

type OrderItem struct {
	ID       uint   `json:"id" gorm:"primaryKey"`
	OrderID  uint   `json:"order_id" gorm:"not null"`
	ItemType string `json:"item_type" gorm:"not null"`
	ItemID   uint   `json:"item_id" gorm:"not null"`
	Jumlah   int    `json:"jumlah" gorm:"not null"`
	Order    Orders `json:"order" gorm:"foreignKey:OrderID"`
}

type BaseRespose struct {
	Status  bool        `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func main() {
	//	loadEnv()
	InitDatabase()

	e := echo.New()
	e.GET("/categories", GetUsersController)
	e.POST("/categories", AddUsersController)
	e.Start(":8000")
}

var DB *gorm.DB

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
	DB.AutoMigrate(
	//&Categories{},
	//&Minuman{},
	//&Makanan{},
	//&Promo{},
	//&Printer{},
	//&Meja{},
	//&Orders{},
	//&OrderItem{},
	)
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

func GetUsersController(c echo.Context) error {
	var users []OrderItem

	// Anda bisa menggunakan Paging jika data terlalu banyak
	// limit := c.QueryParam("limit")
	// offset := c.QueryParam("offset")
	// result := DB.Limit(limit).Offset(offset).Find(&users)

	result := DB.Limit(1).Offset(1).Find(&users)

	if result.Error != nil {
		return c.JSON(http.StatusInternalServerError, BaseRespose{
			Status:  false,
			Message: "Failed to retrieve order details: " + result.Error.Error(),
			Data:    nil,
		})
	}

	return c.JSON(http.StatusOK, BaseRespose{
		Status:  true,
		Message: "Successfully retrieved order details",
		Data:    users,
	})
}

//func loadEnv() {
//	err := godotenv.Load()
//	if err != nil {
//		panic("Failed load env file")
//	}
//}
