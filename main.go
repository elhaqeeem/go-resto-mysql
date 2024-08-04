package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

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
	Nomor int  `json:"nomor" gorm:"not null"`
}

type Orders struct {
	ID      uint      `json:"id" gorm:"primaryKey"`
	MejaID  uint      `json:"meja_id"` // Foreign Key
	Tanggal time.Time `json:"tanggal" gorm:"default:CURRENT_TIMESTAMP"`
	Meja    Meja      `json:"meja" gorm:"foreignKey:MejaID;constraint:OnDelete:CASCADE;"`
	// Relasi satu-ke-banyak
	OrderItems []OrderItems `json:"orderitems" gorm:"foreignKey:OrderID"`
}

type OrderItems struct {
	ID       uint   `json:"id" gorm:"primaryKey"`
	OrderID  uint   `json:"order_id"` // Foreign Key
	ItemType string `json:"itemtype" gorm:"type:enum('Minuman', 'Makanan', 'Promo');not null"`
	ItemID   uint   `json:"items_id"` // Foreign Key
	Jumlah   int    `json:"jumlah" gorm:"not null"`
	Order    Orders `json:"order" gorm:"foreignKey:OrderID;constraint:OnDelete:CASCADE;"`
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
		&Categories{},
		&Minuman{},
		&Makanan{},
		&Promo{},
		&Printer{},
		&Meja{},
		&Orders{},
		&OrderItems{},
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
	var users []OrderItems

	result := DB.Find(&users)

	if result.Error != nil {
		return c.JSON(http.StatusInternalServerError, BaseRespose{
			Status:  false,
			Message: "Failed get data Order Details",
			Data:    nil,
		})
	}

	return c.JSON(http.StatusOK, BaseRespose{
		Status:  true,
		Message: "Success get data Order Details",
		Data:    users,
	})
}

//func loadEnv() {
//	err := godotenv.Load()
//	if err != nil {
//		panic("Failed load env file")
//	}
//}
