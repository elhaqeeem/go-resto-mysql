package main

import (
	"fmt"
	"log"
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

type Orders struct {
	ID      uint      `gorm:"primaryKey"`
	MejaID  uint      // Foreign Key
	Tanggal time.Time `gorm:"default:CURRENT_TIMESTAMP"`
	Meja    Meja      `gorm:"foreignKey:MejaID;constraint:OnDelete:CASCADE;"`
	// Relasi satu-ke-banyak
	OrderItems []OrderItems `gorm:"foreignKey:OrderID"`
}

type OrderItems struct {
	ID       uint   `gorm:"primaryKey"`
	OrderID  uint   // Foreign Key
	ItemType string `gorm:"type:enum('Minuman', 'Makanan', 'Promo');not null"`
	ItemID   uint   // Foreign Key
	Jumlah   int    `gorm:"not null"`
	Order    Orders `gorm:"foreignKey:OrderID;constraint:OnDelete:CASCADE;"`
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
	var meja Meja
	if err := DB.Where("nomor = ?", 1).First(&meja).Error; err != nil {
		log.Fatal("Failed to find table 1:", err)
	}
	order := Orders{MejaID: meja.ID}
	if err := DB.Create(&order).Error; err != nil {
		log.Fatal("Failed to create order:", err)
	}

	// Insert order items
	var orderID = order.ID
	var esBatu, kopiPanas, promoItem, tehManis, mieGoreng Minuman

	if err := DB.Where("nama = ?", "Es Batu").First(&esBatu).Error; err != nil {
		log.Fatal("Failed to find Es Batu:", err)
	}
	if err := DB.Where("nama = ? AND varian = ?", "Kopi", "Panas").First(&kopiPanas).Error; err != nil {
		log.Fatal("Failed to find Kopi Panas:", err)
	}
	if err := DB.Where("nama = ?", "Nasi Goreng + Jeruk Dingin").First(&promoItem).Error; err != nil {
		log.Fatal("Failed to find Promo:", err)
	}
	if err := DB.Where("nama = ? AND varian = ?", "Teh", "Manis").First(&tehManis).Error; err != nil {
		log.Fatal("Failed to find Teh Manis:", err)
	}
	if err := DB.Where("nama = ? AND varian = ?", "Mie", "Goreng").First(&mieGoreng).Error; err != nil {
		log.Fatal("Failed to find Mie Goreng:", err)
	}

	orderItems := []OrderItems{
		{OrderID: orderID, ItemType: "Minuman", ItemID: esBatu.ID, Jumlah: 1},
		{OrderID: orderID, ItemType: "Minuman", ItemID: kopiPanas.ID, Jumlah: 1},
		{OrderID: orderID, ItemType: "Promo", ItemID: promoItem.ID, Jumlah: 2},
		{OrderID: orderID, ItemType: "Minuman", ItemID: tehManis.ID, Jumlah: 1},
		{OrderID: orderID, ItemType: "Makanan", ItemID: mieGoreng.ID, Jumlah: 1},
	}
	if err := DB.Create(&orderItems).Error; err != nil {
		log.Fatal("Failed to insert order items:", err)
	}

	log.Println("Data inserted successfully")
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
