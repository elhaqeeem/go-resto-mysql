package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Categories struct {
	Id        int            `json:"id" gorm:"primaryKey autoIncrement"`
	Nama      string         `json:"nama" gorm:"size:255;uniqueIndex"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
	//Minuman []Minuman `json:"minuman" gorm:"foreignKey:Categories_id"`
	//Makanan []Makanan `json:"makanan" gorm:"foreignKey:Categories_id"`
}

type Minuman struct {
	ID            uint       `json:"id" gorm:"primaryKey"`
	Nama          string     `json:"nama" gorm:"size:50;not null"`
	Varian        string     `json:"varian" gorm:"size:50;uniqueIndex"`
	Harga         float64    `json:"harga" gorm:"not null;type:decimal(10,2)"`
	Categories_id uint       `json:"categories_id"` // Foreign Key
	Categories    Categories `json:"categories" gorm:"foreignKey:Categories_id;constraint:OnDelete:CASCADE;"`
}

type Makanan struct {
	ID            uint       `json:"id" gorm:"primaryKey"`
	Nama          string     `json:"nama" gorm:"size:50;not null"`
	Varian        string     `json:"varian" gorm:"size:50;uniqueIndex"`
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
	loadEnv()
	InitDatabase()

	e := echo.New()
	e.POST("/api/v1/categories", AddCategoriesController)
	e.GET("/api/v1/categories", GetCategoriesController)
	e.PUT("/api/v1/categories", UpdateCategoriesController)
	e.DELETE("/api/v1/categories/:id", SoftDeleteCategoryController)
	e.GET("/api/v1/orderdetails", GetOrderDetailController)
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

// fungsi api categories
func AddCategoriesController(c echo.Context) error {
	var user Categories

	// Bind request data to user
	if err := c.Bind(&user); err != nil {
		return c.JSON(http.StatusBadRequest, BaseRespose{
			Status:  false,
			Message: "Invalid request data",
			Data:    nil,
		})
	}

	// Check if category already exists
	var existingCategory Categories
	if err := DB.Where("nama = ?", user.Nama).First(&existingCategory).Error; err == nil {
		// If no error and a record is found
		formattedMessage := fmt.Sprintf("Categories with nama %s already exists", user.Nama)
		return c.JSON(http.StatusConflict, BaseRespose{
			Status:  false,
			Message: formattedMessage,
			Data:    nil,
		})
	}

	// Create new category
	result := DB.Create(&user)
	if result.Error != nil {
		return c.JSON(http.StatusInternalServerError, BaseRespose{
			Status:  false,
			Message: "Failed to add category",
			Data:    nil,
		})
	}

	return c.JSON(http.StatusCreated, BaseRespose{
		Status:  true,
		Message: "Successfully added category",
		Data:    user,
	})
}

func GetCategoriesController(c echo.Context) error {
	var users []Categories

	result := DB.Order("nama desc").Limit(10).Find(&users)

	if result.Error != nil {
		return c.JSON(http.StatusInternalServerError, BaseRespose{
			Status:  false,
			Message: "Failed to retrieve categories: " + result.Error.Error(),
			Data:    nil,
		})
	}

	return c.JSON(http.StatusOK, BaseRespose{
		Status:  true,
		Message: "Successfully retrieved categories",
		Data:    users,
	})
}

// UpdateCategoriesController updates an existing category
func UpdateCategoriesController(c echo.Context) error {
	var category Categories

	// Bind request data to category
	if err := c.Bind(&category); err != nil {
		return c.JSON(http.StatusBadRequest, BaseRespose{
			Status:  false,
			Message: "Invalid request data",
			Data:    nil,
		})
	}

	// Check if update category already exists

	var existingCategory Categories
	if err := DB.First(&existingCategory, category.Id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.JSON(http.StatusNotFound, BaseRespose{
				Status:  false,
				Message: "Category not found",
				Data:    nil,
			})
		}
		return c.JSON(http.StatusInternalServerError, BaseRespose{
			Status:  false,
			Message: "Failed to check for existing category",
			Data:    nil,
		})
	}

	// Update the existing category
	result := DB.Model(&existingCategory).Updates(category)
	if result.Error != nil {
		return c.JSON(http.StatusInternalServerError, BaseRespose{
			Status:  false,
			Message: "Failed to update category",
			Data:    nil,
		})
	}

	return c.JSON(http.StatusOK, BaseRespose{
		Status:  true,
		Message: "Successfully updated category",
		Data:    existingCategory,
	})
}

func SoftDeleteCategoryController(c echo.Context) error {
	// Extract category ID from request
	id := c.Param("id")

	// Find the category by ID
	var category Categories
	if err := DB.First(&category, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.JSON(http.StatusNotFound, BaseRespose{
				Status:  false,
				Message: "Category not found",
				Data:    nil,
			})
		}
		return c.JSON(http.StatusInternalServerError, BaseRespose{
			Status:  false,
			Message: "Failed to find category",
			Data:    nil,
		})
	}

	// Perform soft delete
	if err := DB.Delete(&category).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, BaseRespose{
			Status:  false,
			Message: "Failed to delete category",
			Data:    nil,
		})
	}

	return c.JSON(http.StatusOK, BaseRespose{
		Status:  true,
		Message: "Successfully deleted category",
		Data:    nil,
	})
}

func GetOrderDetailController(c echo.Context) error {
	var users []OrderItem

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

func loadEnv() {
	err := godotenv.Load()
	if err != nil {
		panic("Failed load env file")
	}
}
