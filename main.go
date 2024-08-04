package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Categories struct {
	Id      int       `json:"id" gorm:"primaryKey autoIncrement"`
	Nama    string    `json:"nama"`
	Minuman []Minuman `gorm:"foreignKey:Categories_id"`
	Makanan []Makanan `gorm:"foreignKey:Categories_id"`
}

type Minuman struct {
	ID            uint       `gorm:"primaryKey"`
	Nama          string     `gorm:"size:50;not null"`
	Varian        string     `gorm:"size:50"`
	Harga         float64    `gorm:"not null;type:decimal(10,2)"`
	Categories_id uint       // Foreign Key
	Categories    Categories `gorm:"foreignKey:Categories_id;constraint:OnDelete:CASCADE;"`
}

type Makanan struct {
	ID            uint       `gorm:"primaryKey"`
	Nama          string     `gorm:"size:50;not null"`
	Varian        string     `gorm:"size:50"`
	Harga         float64    `gorm:"not null;type:decimal(10,2)"`
	Categories_id uint       // Foreign Key
	Categories    Categories `gorm:"foreignKey:Categories_id;constraint:OnDelete:CASCADE;"`
}

type Promo struct {
	ID    uint    `gorm:"primaryKey"`
	Nama  string  `gorm:"size:100;not null"`
	Harga float64 `gorm:"not null;type:decimal(10,2)"`
}

type Printer struct {
	ID   string `gorm:"primaryKey;size:1"`
	Nama string `gorm:"size:50;not null"`
}

type Meja struct {
	ID    uint `gorm:"primaryKey"`
	Nomor int  `gorm:"not null"`
}

type BaseRespose struct {
	Status  bool        `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func main() {
	//	loadEnv()
	InitDatabase()
	categories := []Categories{
		{Nama: "Minuman"},
		{Nama: "Makanan"},
		{Nama: "Promo"},
	}
	DB.Create(&categories)

	// Retrieve category IDs
	var minumanCategory, makananCategory Categories
	DB.Where("nama = ?", "Minuman").First(&minumanCategory)
	DB.Where("nama = ?", "Makanan").First(&makananCategory)

	// Insert drinks
	drinks := []Minuman{
		{Nama: "Jeruk", Varian: "Dingin", Harga: 12000, Categories_id: uint(minumanCategory.Id)},
		{Nama: "Jeruk", Varian: "Panas", Harga: 10000, Categories_id: uint(minumanCategory.Id)},
		{Nama: "Teh", Varian: "Manis", Harga: 8000, Categories_id: uint(minumanCategory.Id)},
		{Nama: "Teh", Varian: "Tawar", Harga: 5000, Categories_id: uint(minumanCategory.Id)},
		{Nama: "Kopi", Varian: "Dingin", Harga: 8000, Categories_id: uint(minumanCategory.Id)},
		{Nama: "Kopi", Varian: "Panas", Harga: 6000, Categories_id: uint(minumanCategory.Id)},
		{Nama: "Es Batu", Varian: "", Harga: 2000, Categories_id: uint(minumanCategory.Id)},
	}
	DB.Create(&drinks)

	// Insert foods
	foods := []Makanan{
		{Nama: "Mie", Varian: "Goreng", Harga: 15000, Categories_id: uint(makananCategory.Id)},
		{Nama: "Mie", Varian: "Kuah", Harga: 15000, Categories_id: uint(makananCategory.Id)},
		{Nama: "Nasi Goreng", Varian: "", Harga: 15000, Categories_id: uint(makananCategory.Id)},
	}
	DB.Create(&foods)

	// Insert promo
	promo := Promo{Nama: "Nasi Goreng + Jeruk Dingin", Harga: 23000}
	DB.Create(&promo)

	// Insert printers
	printers := []Printer{
		{ID: "A", Nama: "Printer Kasir"},
		{ID: "B", Nama: "Printer Dapur (Makanan)"},
		{ID: "C", Nama: "Printer Bar (Minuman)"},
	}
	DB.Create(&printers)

	// Insert tables
	tables := []Meja{
		{Nomor: 1},
		{Nomor: 2},
		{Nomor: 3},
	}
	DB.Create(&tables)

	log.Println("Records created successfully")

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
		&Categories{},
		&Makanan{},
		&Minuman{},
		&Promo{},
		&Printer{},
		&Meja{},
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

//func loadEnv() {
//	err := godotenv.Load()
//	if err != nil {
//		panic("Failed load env file")
//	}
//}
