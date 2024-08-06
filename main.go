package main

import (
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Categories struct {
	Id        int            `json:"id" gorm:"primaryKey autoIncrement"`
	Nama      string         `json:"nama" gorm:"size:255;uniqueIndex"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
	//Minuman   []Minuman      `json:"minuman" gorm:"foreignKey:Categories_id"`
	//Makanan   []Makanan      `json:"makanan" gorm:"foreignKey:Categories_id"`
}

type Makanan struct {
	ID            uint           `json:"id" gorm:"primaryKey"`
	Nama          string         `json:"nama" gorm:"size:50;not null"`
	Varian        string         `json:"varian" gorm:"size:50;uniqueIndex"`
	Harga         float64        `json:"harga" gorm:"not null;type:decimal(10,2)"`
	Categories_id uint           `json:"categories_id"` // Foreign Key
	Categories    Categories     `json:"categories" gorm:"foreignKey:Categories_id;constraint:OnDelete:CASCADE;"`
	DeletedAt     gorm.DeletedAt `gorm:"index"`
}

type Minuman struct {
	ID            uint           `json:"id" gorm:"primaryKey"`
	Nama          string         `json:"nama" gorm:"size:50;not null"`
	Varian        string         `json:"varian" gorm:"size:50;uniqueIndex"`
	Harga         float64        `json:"harga" gorm:"not null;type:decimal(10,2)"`
	Categories_id uint           `json:"categories_id"` // Foreign Key
	Categories    Categories     `json:"categories" gorm:"foreignKey:Categories_id;constraint:OnDelete:CASCADE;"`
	DeletedAt     gorm.DeletedAt `gorm:"index"`
}

type Promo struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	Nama      string         `json:"nama" gorm:"size:100;;uniqueIndex"`
	Harga     float64        `json:"harga" gorm:"not null;type:decimal(10,2)"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

type Printer struct {
	ID        string         `json:"id" gorm:"primaryKey;size:1"`
	Nama      string         `json:"nama" gorm:"size:50;not null"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

type Meja struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	Nomor     int            `json:"nomor" gorm:"unique;not null"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

type Orders struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	MejaID    uint           `json:"meja_id" gorm:"not null"`
	Meja      Meja           `json:"meja" gorm:"foreignKey:MejaID"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

type OrderItem struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	OrderID   uint           `json:"order_id" gorm:"not null"`
	ItemType  string         `json:"item_type" gorm:"not null"`
	ItemID    uint           `json:"item_id" gorm:"not null"`
	Jumlah    int            `json:"jumlah" gorm:"not null"`
	Order     Orders         `json:"order" gorm:"foreignKey:OrderID"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

type BaseRespose struct {
	Status  bool        `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func main() {
	//loadEnv()
	InitDatabase()

	e := echo.New()
	//route api Categories
	e.POST("/api/v1/categories", AddCategoriesController)
	e.GET("/api/v1/categories", GetCategoriesController)
	e.PUT("/api/v1/categories", UpdateCategoriesController)
	e.DELETE("/api/v1/categories/:id", SoftDeleteCategoryController)
	e.PUT("/api/v1/categories/restore/:id", RestoreCategoryController)
	e.DELETE("/api/v1/categories/:id", DeleteCategoryController)
	//route api Makanan
	e.POST("/api/v1/food", AddMakananController)
	e.GET("/api/v1/food", GetMakananController)
	e.PUT("/api/v1/food", UpdateMakananController)
	e.DELETE("/api/v1/food/:id", SoftDeleteMakananController)
	e.PUT("/api/v1/food/restore/:id", RestoreMakananController)
	e.DELETE("/api/v1/food/:id", DeleteMakananController)
	//route api Minuman
	e.POST("/api/v1/drink", AddMinumanController)
	e.GET("/api/v1/drink", GetMinumanController)
	e.PUT("/api/v1/drink", UpdateMinumanController)
	e.DELETE("/api/v1/drink/:id", SoftDeleteMinumanController)
	e.PUT("/api/v1/drink/restore/:id", RestoreMinumanController)
	e.DELETE("/api/v1/drink/:id", DeleteMinumanController)
	//route api Promo
	e.POST("/api/v1/discount", CreatePromo)
	e.GET("/api/v1/discount", GetPromoByID)
	e.PUT("/api/v1/discount", UpdatePromo)
	e.DELETE("/api/v1/discount/:id", DeletePromo)
	//route api Meja
	//route api Printer

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

// controller categories
// add
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

// read
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

// soft delete categories
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

// restore data categories
func RestoreCategoryController(c echo.Context) error {
	// Extract category ID from request
	id := c.Param("id")

	// Find the category by ID (including soft-deleted records)
	var category Categories
	if err := DB.Unscoped().First(&category, id).Error; err != nil {
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

	// Check if the category is already active
	if category.DeletedAt.Time.IsZero() {
		return c.JSON(http.StatusConflict, BaseRespose{
			Status:  false,
			Message: "Category is not deleted",
			Data:    nil,
		})
	}

	// Restore the category by setting DeletedAt to nil
	if err := DB.Model(&category).Update("DeletedAt", nil).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, BaseRespose{
			Status:  false,
			Message: "Failed to restore category",
			Data:    nil,
		})
	}

	return c.JSON(http.StatusOK, BaseRespose{
		Status:  true,
		Message: "Successfully restored category",
		Data:    category,
	})
}

// delete permanen data categories
func DeleteCategoryController(c echo.Context) error {
	id := c.Param("id")

	// Find the Makanan entry by ID
	var categories Categories
	if err := DB.First(&categories, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.JSON(http.StatusNotFound, BaseRespose{
				Status:  false,
				Message: "Categories not found",
				Data:    nil,
			})
		}
		return c.JSON(http.StatusInternalServerError, BaseRespose{
			Status:  false,
			Message: "Failed to find Categories",
			Data:    nil,
		})
	}

	// Delete the Makanan entry
	if err := DB.Delete(&categories).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, BaseRespose{
			Status:  false,
			Message: "Failed to delete Categories",
			Data:    nil,
		})
	}

	return c.JSON(http.StatusOK, BaseRespose{
		Status:  true,
		Message: "Successfully deleted Categories",
		Data:    nil,
	})
}

// controller makanan
// add makanan
func AddMakananController(c echo.Context) error {
	var request struct {
		Nama         string  `json:"nama"`
		Varian       string  `json:"varian"`
		Harga        float64 `json:"harga"`
		CategoriesID uint    `json:"categories_id"`
	}

	// Bind JSON request to struct
	if err := c.Bind(&request); err != nil {
		return c.JSON(http.StatusBadRequest, BaseRespose{
			Status:  false,
			Message: "Invalid request data",
			Data:    nil,
		})
	}

	// Create new Makanan entry
	makanan := Makanan{
		Nama:          request.Nama,
		Varian:        request.Varian,
		Harga:         request.Harga,
		Categories_id: request.CategoriesID,
	}

	// Save to database
	if err := DB.Create(&makanan).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, BaseRespose{
			Status:  false,
			Message: "Failed to add makanan",
			Data:    nil,
		})
	}

	return c.JSON(http.StatusCreated, BaseRespose{
		Status:  true,
		Message: "Successfully added makanan",
		Data:    makanan,
	})
}

// read
func GetMakananController(c echo.Context) error {
	id := c.Param("id")

	// Find the Makanan entry by ID
	var makanan Makanan
	if err := DB.Preload("Categories").First(&makanan, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.JSON(http.StatusNotFound, BaseRespose{
				Status:  false,
				Message: "Makanan not found",
				Data:    nil,
			})
		}
		return c.JSON(http.StatusInternalServerError, BaseRespose{
			Status:  false,
			Message: "Failed to retrieve makanan",
			Data:    nil,
		})
	}

	return c.JSON(http.StatusOK, BaseRespose{
		Status:  true,
		Message: "Successfully retrieved makanan",
		Data:    makanan,
	})
}

// update
func UpdateMakananController(c echo.Context) error {
	id := c.Param("id")
	var request struct {
		Nama         string  `json:"nama"`
		Varian       string  `json:"varian"`
		Harga        float64 `json:"harga"`
		CategoriesID uint    `json:"categories_id"`
	}

	// Bind JSON request to struct
	if err := c.Bind(&request); err != nil {
		return c.JSON(http.StatusBadRequest, BaseRespose{
			Status:  false,
			Message: "Invalid request data",
			Data:    nil,
		})
	}

	// Find the existing Makanan entry
	var makanan Makanan
	if err := DB.First(&makanan, id).Error; err != nil {
		return c.JSON(http.StatusNotFound, BaseRespose{
			Status:  false,
			Message: "Makanan not found",
			Data:    nil,
		})
	}

	// Update fields
	makanan.Nama = request.Nama
	makanan.Varian = request.Varian
	makanan.Harga = request.Harga
	makanan.Categories_id = request.CategoriesID

	// Save updates to database
	if err := DB.Save(&makanan).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, BaseRespose{
			Status:  false,
			Message: "Failed to update makanan",
			Data:    nil,
		})
	}

	return c.JSON(http.StatusOK, BaseRespose{
		Status:  true,
		Message: "Successfully updated makanan",
		Data:    makanan,
	})
}

// softdelele makanan
func SoftDeleteMakananController(c echo.Context) error {
	// Extract category ID from request
	id := c.Param("id")

	// Find the category by ID
	var makanan Makanan
	if err := DB.First(&makanan, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.JSON(http.StatusNotFound, BaseRespose{
				Status:  false,
				Message: "Food not found",
				Data:    nil,
			})
		}
		return c.JSON(http.StatusInternalServerError, BaseRespose{
			Status:  false,
			Message: "Failed to find Food",
			Data:    nil,
		})
	}

	// Perform soft delete
	if err := DB.Delete(&makanan).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, BaseRespose{
			Status:  false,
			Message: "Failed to delete food",
			Data:    nil,
		})
	}

	return c.JSON(http.StatusOK, BaseRespose{
		Status:  true,
		Message: "Successfully deleted food",
		Data:    nil,
	})
}

// restore makanan
func RestoreMakananController(c echo.Context) error {
	// Extract category ID from request
	id := c.Param("id")

	// Find the category by ID (including soft-deleted records)
	var makanan Makanan
	if err := DB.Unscoped().First(&makanan, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.JSON(http.StatusNotFound, BaseRespose{
				Status:  false,
				Message: "Food not found",
				Data:    nil,
			})
		}
		return c.JSON(http.StatusInternalServerError, BaseRespose{
			Status:  false,
			Message: "Failed to find food",
			Data:    nil,
		})
	}

	// Check if the category is already active
	if makanan.DeletedAt.Time.IsZero() {
		return c.JSON(http.StatusConflict, BaseRespose{
			Status:  false,
			Message: "Food is not deleted",
			Data:    nil,
		})
	}

	// Restore the category by setting DeletedAt to nil
	if err := DB.Model(&makanan).Update("DeletedAt", nil).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, BaseRespose{
			Status:  false,
			Message: "Failed to restore category",
			Data:    nil,
		})
	}

	return c.JSON(http.StatusOK, BaseRespose{
		Status:  true,
		Message: "Successfully restored category",
		Data:    makanan,
	})
}

// delete  makanan
func DeleteMakananController(c echo.Context) error {
	id := c.Param("id")

	// Find the Makanan entry by ID
	var makanan Makanan
	if err := DB.First(&makanan, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.JSON(http.StatusNotFound, BaseRespose{
				Status:  false,
				Message: "Makanan not found",
				Data:    nil,
			})
		}
		return c.JSON(http.StatusInternalServerError, BaseRespose{
			Status:  false,
			Message: "Failed to find makanan",
			Data:    nil,
		})
	}

	// Delete the Makanan entry
	if err := DB.Delete(&makanan).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, BaseRespose{
			Status:  false,
			Message: "Failed to delete makanan",
			Data:    nil,
		})
	}

	return c.JSON(http.StatusOK, BaseRespose{
		Status:  true,
		Message: "Successfully deleted makanan",
		Data:    nil,
	})
}

// controller minuman
// add minuman
func AddMinumanController(c echo.Context) error {
	var request struct {
		Nama         string  `json:"nama"`
		Varian       string  `json:"varian"`
		Harga        float64 `json:"harga"`
		CategoriesID uint    `json:"categories_id"`
	}

	// Bind JSON request to struct
	if err := c.Bind(&request); err != nil {
		return c.JSON(http.StatusBadRequest, BaseRespose{
			Status:  false,
			Message: "Invalid request data",
			Data:    nil,
		})
	}

	// Create new Makanan entry
	minuman := Minuman{
		Nama:          request.Nama,
		Varian:        request.Varian,
		Harga:         request.Harga,
		Categories_id: request.CategoriesID,
	}

	// Save to database
	if err := DB.Create(&minuman).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, BaseRespose{
			Status:  false,
			Message: "Failed to add Drink",
			Data:    nil,
		})
	}

	return c.JSON(http.StatusCreated, BaseRespose{
		Status:  true,
		Message: "Successfully added Drink",
		Data:    minuman,
	})
}

// read minuman
func GetMinumanController(c echo.Context) error {
	id := c.Param("id")

	// Find the Makanan entry by ID
	var minuman Minuman
	if err := DB.Preload("minumen").First(&minuman, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.JSON(http.StatusNotFound, BaseRespose{
				Status:  false,
				Message: "minuman not found",
				Data:    nil,
			})
		}
		return c.JSON(http.StatusInternalServerError, BaseRespose{
			Status:  false,
			Message: "Failed to retrieve minuman",
			Data:    nil,
		})
	}

	return c.JSON(http.StatusOK, BaseRespose{
		Status:  true,
		Message: "Successfully retrieved minuman",
		Data:    minuman,
	})
}

// update minuman
func UpdateMinumanController(c echo.Context) error {
	id := c.Param("id")
	var request struct {
		Nama         string  `json:"nama"`
		Varian       string  `json:"varian"`
		Harga        float64 `json:"harga"`
		CategoriesID uint    `json:"categories_id"`
	}

	// Bind JSON request to struct
	if err := c.Bind(&request); err != nil {
		return c.JSON(http.StatusBadRequest, BaseRespose{
			Status:  false,
			Message: "Invalid request data",
			Data:    nil,
		})
	}

	// Find the existing Makanan entry
	var minuman Minuman
	if err := DB.First(&minuman, id).Error; err != nil {
		return c.JSON(http.StatusNotFound, BaseRespose{
			Status:  false,
			Message: "Makanan not found",
			Data:    nil,
		})
	}

	// Update fields
	minuman.Nama = request.Nama
	minuman.Varian = request.Varian
	minuman.Harga = request.Harga
	minuman.Categories_id = request.CategoriesID

	// Save updates to database
	if err := DB.Save(&minuman).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, BaseRespose{
			Status:  false,
			Message: "Failed to update minuman",
			Data:    nil,
		})
	}

	return c.JSON(http.StatusOK, BaseRespose{
		Status:  true,
		Message: "Successfully updated minuman",
		Data:    minuman,
	})
}

// softdelete minuman
func SoftDeleteMinumanController(c echo.Context) error {
	// Extract category ID from request
	id := c.Param("id")

	// Find the category by ID
	var minuman Minuman
	if err := DB.First(&minuman, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.JSON(http.StatusNotFound, BaseRespose{
				Status:  false,
				Message: "drink not found",
				Data:    nil,
			})
		}
		return c.JSON(http.StatusInternalServerError, BaseRespose{
			Status:  false,
			Message: "Failed to find drink",
			Data:    nil,
		})
	}

	// Perform soft delete
	if err := DB.Delete(&minuman).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, BaseRespose{
			Status:  false,
			Message: "Failed to delete drink",
			Data:    nil,
		})
	}

	return c.JSON(http.StatusOK, BaseRespose{
		Status:  true,
		Message: "Successfully deleted drink",
		Data:    nil,
	})
}

// restore minuman
func RestoreMinumanController(c echo.Context) error {
	// Extract category ID from request
	id := c.Param("id")

	// Find the category by ID (including soft-deleted records)
	var minuman Minuman
	if err := DB.Unscoped().First(&minuman, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.JSON(http.StatusNotFound, BaseRespose{
				Status:  false,
				Message: "drink not found",
				Data:    nil,
			})
		}
		return c.JSON(http.StatusInternalServerError, BaseRespose{
			Status:  false,
			Message: "Failed to find drink",
			Data:    nil,
		})
	}

	// Check if the category is already active
	if minuman.DeletedAt.Time.IsZero() {
		return c.JSON(http.StatusConflict, BaseRespose{
			Status:  false,
			Message: "drink is not deleted",
			Data:    nil,
		})
	}

	// Restore the category by setting DeletedAt to nil
	if err := DB.Model(&minuman).Update("DeletedAt", nil).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, BaseRespose{
			Status:  false,
			Message: "Failed to restore category",
			Data:    nil,
		})
	}

	return c.JSON(http.StatusOK, BaseRespose{
		Status:  true,
		Message: "Successfully restored category",
		Data:    minuman,
	})
}

// delete minuman
func DeleteMinumanController(c echo.Context) error {
	id := c.Param("id")

	// Find the Makanan entry by ID
	var minuman Minuman
	if err := DB.First(&minuman, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.JSON(http.StatusNotFound, BaseRespose{
				Status:  false,
				Message: "Drink not found",
				Data:    nil,
			})
		}
		return c.JSON(http.StatusInternalServerError, BaseRespose{
			Status:  false,
			Message: "Failed to find drink",
			Data:    nil,
		})
	}

	// Delete the Makanan entry
	if err := DB.Delete(&minuman).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, BaseRespose{
			Status:  false,
			Message: "Failed to delete minuman",
			Data:    nil,
		})
	}

	return c.JSON(http.StatusOK, BaseRespose{
		Status:  true,
		Message: "Successfully deleted minuman",
		Data:    nil,
	})
}

// controller promo
// add
func CreatePromo(c echo.Context) error {
	var promo Promo

	if err := c.Bind(&promo); err != nil {
		return c.JSON(http.StatusBadRequest, BaseRespose{
			Status:  false,
			Message: "Invalid request data",
			Data:    nil,
		})
	}

	if promo.Nama == "" || promo.Harga <= 0 {
		return c.JSON(http.StatusBadRequest, BaseRespose{
			Status:  false,
			Message: "Invalid promo data",
			Data:    nil,
		})
	}

	result := DB.Create(&promo)
	if result.Error != nil {
		return c.JSON(http.StatusInternalServerError, BaseRespose{
			Status:  false,
			Message: "Failed to create promo",
			Data:    nil,
		})
	}

	return c.JSON(http.StatusCreated, BaseRespose{
		Status:  true,
		Message: "Promo created successfully",
		Data:    promo,
	})
}

// read
func GetPromoByID(c echo.Context) error {
	id := c.Param("id")
	var promo Promo

	if err := DB.First(&promo, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.JSON(http.StatusNotFound, BaseRespose{
				Status:  false,
				Message: "Promo not found",
				Data:    nil,
			})
		}
		return c.JSON(http.StatusInternalServerError, BaseRespose{
			Status:  false,
			Message: "Error retrieving promo",
			Data:    nil,
		})
	}

	return c.JSON(http.StatusOK, BaseRespose{
		Status:  true,
		Message: "Promo retrieved successfully",
		Data:    promo,
	})
}

// update
func UpdatePromo(c echo.Context) error {
	id := c.Param("id")
	var updatedPromo Promo

	if err := c.Bind(&updatedPromo); err != nil {
		return c.JSON(http.StatusBadRequest, BaseRespose{
			Status:  false,
			Message: "Invalid request data",
			Data:    nil,
		})
	}

	var existingPromo Promo
	if err := DB.First(&existingPromo, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.JSON(http.StatusNotFound, BaseRespose{
				Status:  false,
				Message: "Promo not found",
				Data:    nil,
			})
		}
		return c.JSON(http.StatusInternalServerError, BaseRespose{
			Status:  false,
			Message: "Error retrieving promo",
			Data:    nil,
		})
	}

	// Update fields
	existingPromo.Nama = updatedPromo.Nama
	existingPromo.Harga = updatedPromo.Harga

	result := DB.Save(&existingPromo)
	if result.Error != nil {
		return c.JSON(http.StatusInternalServerError, BaseRespose{
			Status:  false,
			Message: "Failed to update promo",
			Data:    nil,
		})
	}

	return c.JSON(http.StatusOK, BaseRespose{
		Status:  true,
		Message: "Promo updated successfully",
		Data:    existingPromo,
	})
}

// softdelele
// restore
// delete
func DeletePromo(c echo.Context) error {
	id := c.Param("id")
	result := DB.Delete(&Promo{}, id)

	if result.Error != nil {
		return c.JSON(http.StatusInternalServerError, BaseRespose{
			Status:  false,
			Message: "Failed to delete promo",
			Data:    nil,
		})
	}

	if result.RowsAffected == 0 {
		return c.JSON(http.StatusNotFound, BaseRespose{
			Status:  false,
			Message: "Promo not found",
			Data:    nil,
		})
	}

	return c.JSON(http.StatusOK, BaseRespose{
		Status:  true,
		Message: "Promo deleted successfully",
		Data:    nil,
	})
}

//func loadEnv() {
//	err := godotenv.Load()
//	if err != nil {
//		panic("Failed load env file")
//	}
//}
