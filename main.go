package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-redis/redis/v8"
	"github.com/labstack/echo/v4"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// Product represents a product in the database
type Product struct {
	gorm.Model
	ID       uint    `gorm:"primaryKey"`
	Category string  `gorm:"not null"`
	Name     string  `gorm:"not null"`
	Varian   string  `gorm:"not null"`
	Price    float64 `gorm:"not null;type:decimal(10,2)"`
}

type UpdateOrderRequest struct {
	TableNumber int `json:"table_number"`
	Status      int `json:"status"`
}

// Order represents an order with its items
type Order struct {
	gorm.Model
	ID          uint `gorm:"primaryKey"`
	TableNumber int  // Use int here
	Status      int
	Items       []OrderItem
}

// OrderItem represents an item in an order
type OrderItem struct {
	ID        uint `gorm:"primaryKey"`
	OrderID   uint
	ProductID uint
	Quantity  int
	Product   Product `gorm:"foreignKey:ProductID;references:ID"`
}

// OrderPrinter represents printers assigned to an order
type OrderPrinter struct {
	gorm.Model
	OrderID   uint    `gorm:"not null"`
	PrinterID string  `gorm:"size:1;not null"`
	Printer   Printer `gorm:"foreignKey:PrinterID;references:ID"`
}

// Printer represents a printer in the system
type Printer struct {
	gorm.Model
	ID   string `gorm:"size:1;primaryKey"`
	Name string `gorm:"size:50;uniqueIndex"`
}

// Promo represents a promotional discount
type Promo struct {
	ID         uint           `gorm:"primaryKey"`
	Nama       string         `gorm:"size:100;uniqueIndex"`
	Harga      float64        `gorm:"not null;type:decimal(10,2)"`
	ProductIDs []uint         `gorm:"-"` // Excluded from DB schema
	DeletedAt  gorm.DeletedAt `gorm:"index"`
}

// Meja represents a table in the restaurant
type Meja struct {
	gorm.Model
	ID   uint   `gorm:"primaryKey"`
	Nama string `gorm:"size:50;uniqueIndex"` // Ensures unique Nama

}

// CreateOrderRequest is used for creating a new order
type CreateOrderRequest struct {
	TableNumber int                `json:"table_number"`
	Items       []OrderItemRequest `json:"items"`
}

// OrderItemRequest represents each item in the order request
type OrderItemRequest struct {
	ProductID uint `json:"product_id"`
	Quantity  int  `json:"quantity"`
}

type CreateOrderResponse struct {
	Status  bool      `json:"status"`
	Message string    `json:"message"`
	Data    OrderData `json:"data"`
}

type OrderData struct {
	DebugInfo   []string `json:"debug_info"`
	OrderID     uint     `json:"order_id"`
	Printers    []string `json:"printers"`
	TableNumber int      `json:"table_number"`
}

type BaseResponse struct {
	Status  bool        `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func main() {
	//loadEnv()
	InitDatabase()
	//InitRedis()
	e := echo.New()

	//route api Promo
	e.POST("/api/v1/discount", AddPromoController)
	e.GET("/api/v1/discount", GetPromoByIDController)
	e.PUT("/api/v1/discount", UpdatePromoController)
	e.DELETE("/api/v1/discount/:id", DeletePromoController)
	//route api Meja
	e.POST("/api/v1/table", AddMejaController)
	e.GET("/api/v1/table", GetMejasController)
	e.PUT("/api/v1/table/:id", UpdateMejaController)
	e.DELETE("/api/v1/table/:id", SoftDeleteMejaController)
	e.PUT("/api/v1/table/restore/:id", RestoreMejaController)
	e.DELETE("/api/v1/table/hard-delete/:id", DeleteMejaController)

	//route api Printer
	e.POST("/api/v1/printers", CreatePrinterController)
	e.GET("api/v1/printers", GetPrintersController)
	e.PUT("/api/v1/printers", UpdatePrinterController)
	e.DELETE("/api/v1/printers/:id", SoftDeletePrinterController)
	e.PUT("/api/v1/printers/:id/restore", RestorePrinterController)
	e.DELETE("/api/v1/printers/hard-delete/:id", DeletePrinterController)
	//post menu
	e.POST("/api/v1/product", CreateProductController)
	e.GET("api/v1/product", GetProductsController)
	e.PUT("/api/v1/product", UpdateProductController)
	e.DELETE("/api/v1/products/:id/soft-delete", SoftDeleteProductController)
	e.PUT("/api/v1/product/:id/restore", RestoreProductController)
	e.DELETE("/api/v1/product/hard-delete/:id", DeleteProductController)
	//post order
	e.POST("/api/v1/neworder", CreateOrder)
	e.GET("api/v1/neworder", GetOrderController)
	e.PUT("/api/v1/neworder", UpdateOrderController)
	e.DELETE("/api/v1/neworder/:id", SoftDeleteOrderController)
	e.PUT("/api/v1/neworder/restore/:id", RestoreOrderController)
	e.DELETE("/api/v1/neworder/hard-delete/:id", DeleteOrderController)
	//route api Get bill
	e.GET("/api/v1/bill/:table_number", GetBill)
	e.Start(":8000")
}

var redisClient *redis.Client
var ctx = context.Background()
var DB *gorm.DB

func InitDatabase() {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"))

	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("Failed to initialize database: " + err.Error())
	}
	Migration()
}
func InitRedis() {
	redisClient = redis.NewClient(&redis.Options{
		Addr:     "redis:6379", // Redis service name and port
		Password: "",           // No password set
		DB:       0,            // Default DB
	})

	_, err := redisClient.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}
	fmt.Println("Connected to Redis")
}
func Migration() {
	DB.AutoMigrate(
	//&Promo{},
	//&Printer{},
	//&Meja{},
	//&Product{},
	//&Order{},
	//&OrderItem{},
	//&OrderItemRequest{},
	//&OrderPrinter{},
	//&CreateOrderRequest{},
	//&CreateOrderResponse{},
	//&OrderData{},
	)

}

// controller promo
// add
func AddPromoController(c echo.Context) error {
	var promo Promo

	// Bind request data to promo first
	if err := c.Bind(&promo); err != nil {
		return c.JSON(http.StatusBadRequest, BaseResponse{
			Status:  false,
			Message: "Invalid request data",
			Data:    nil,
		})
	}

	// Validate promo data
	if promo.Nama == "" || promo.Harga <= 0 {
		return c.JSON(http.StatusBadRequest, BaseResponse{
			Status:  false,
			Message: "Invalid promo data",
			Data:    nil,
		})
	}

	// Check if promo with the same name already exists
	var existingPromo Promo
	if err := DB.Where("nama = ?", promo.Nama).First(&existingPromo).Error; err == nil {
		// If no error and a record is found
		formattedMessage := fmt.Sprintf("Promo with nama %s already exists", promo.Nama)
		return c.JSON(http.StatusConflict, BaseResponse{
			Status:  false,
			Message: formattedMessage,
			Data:    nil,
		})
	}

	// Create new promo
	result := DB.Create(&promo)
	if result.Error != nil {
		return c.JSON(http.StatusInternalServerError, BaseResponse{
			Status:  false,
			Message: "Failed to create promo",
			Data:    nil,
		})
	}

	return c.JSON(http.StatusCreated, BaseResponse{
		Status:  true,
		Message: "Promo created successfully",
		Data:    promo,
	})
}

// read
func GetPromoByIDController(c echo.Context) error {
	id := c.Param("id")
	var promo Promo

	if err := DB.First(&promo, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.JSON(http.StatusNotFound, BaseResponse{
				Status:  false,
				Message: "Promo not found",
				Data:    nil,
			})
		}
		return c.JSON(http.StatusInternalServerError, BaseResponse{
			Status:  false,
			Message: "Error retrieving promo",
			Data:    nil,
		})
	}

	return c.JSON(http.StatusOK, BaseResponse{
		Status:  true,
		Message: "Promo retrieved successfully",
		Data:    promo,
	})
}

// update
func UpdatePromoController(c echo.Context) error {
	id := c.Param("id")
	var updatedPromo Promo

	if err := c.Bind(&updatedPromo); err != nil {
		return c.JSON(http.StatusBadRequest, BaseResponse{
			Status:  false,
			Message: "Invalid request data",
			Data:    nil,
		})
	}

	var existingPromo Printer
	if err := DB.First(&existingPromo, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.JSON(http.StatusNotFound, BaseResponse{
				Status:  false,
				Message: "Promo not found",
				Data:    nil,
			})
		}
		return c.JSON(http.StatusInternalServerError, BaseResponse{
			Status:  false,
			Message: "Error retrieving promo",
			Data:    nil,
		})
	}

	// Update fields
	existingPromo.Name = updatedPromo.Nama

	result := DB.Save(&existingPromo)
	if result.Error != nil {
		return c.JSON(http.StatusInternalServerError, BaseResponse{
			Status:  false,
			Message: "Failed to update promo",
			Data:    nil,
		})
	}

	return c.JSON(http.StatusOK, BaseResponse{
		Status:  true,
		Message: "Promo updated successfully",
		Data:    existingPromo,
	})
}

// softdelele
func SoftDeletePromoController(c echo.Context) error {
	// Extract category ID from request
	id := c.Param("id")

	// Find the category by ID
	var promo Printer
	if err := DB.First(&promo, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.JSON(http.StatusNotFound, BaseResponse{
				Status:  false,
				Message: "Promo not found",
				Data:    nil,
			})
		}
		return c.JSON(http.StatusInternalServerError, BaseResponse{
			Status:  false,
			Message: "Failed to find Promo",
			Data:    nil,
		})
	}

	// Perform soft delete
	if err := DB.Delete(&promo).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, BaseResponse{
			Status:  false,
			Message: "Failed to delete category",
			Data:    nil,
		})
	}

	return c.JSON(http.StatusOK, BaseResponse{
		Status:  true,
		Message: "Successfully deleted category",
		Data:    nil,
	})
}

// restore
func RestorePromoController(c echo.Context) error {
	// Extract category ID from request
	id := c.Param("id")

	// Find the category by ID (including soft-deleted records)
	var promo Printer
	if err := DB.Unscoped().First(&promo, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.JSON(http.StatusNotFound, BaseResponse{
				Status:  false,
				Message: "Promo not found",
				Data:    nil,
			})
		}
		return c.JSON(http.StatusInternalServerError, BaseResponse{
			Status:  false,
			Message: "Failed to find promo",
			Data:    nil,
		})
	}

	// Check if the category is already active
	if promo.DeletedAt.Time.IsZero() {
		return c.JSON(http.StatusConflict, BaseResponse{
			Status:  false,
			Message: "promo is not deleted",
			Data:    nil,
		})
	}

	// Restore the category by setting DeletedAt to nil
	if err := DB.Model(&promo).Update("DeletedAt", nil).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, BaseResponse{
			Status:  false,
			Message: "Failed to restore Promo",
			Data:    nil,
		})
	}

	return c.JSON(http.StatusOK, BaseResponse{
		Status:  true,
		Message: "Successfully restored category",
		Data:    promo,
	})
}

// delete
func DeletePromoController(c echo.Context) error {
	id := c.Param("id")
	result := DB.Delete(&Printer{}, id)

	if result.Error != nil {
		return c.JSON(http.StatusInternalServerError, BaseResponse{
			Status:  false,
			Message: "Failed to delete promo",
			Data:    nil,
		})
	}

	if result.RowsAffected == 0 {
		return c.JSON(http.StatusNotFound, BaseResponse{
			Status:  false,
			Message: "Promo not found",
			Data:    nil,
		})
	}

	return c.JSON(http.StatusOK, BaseResponse{
		Status:  true,
		Message: "Promo deleted successfully",
		Data:    nil,
	})
}

// CreatePrinter creates a new printer record
func CreatePrinterController(c echo.Context) error {
	var printer Printer

	// Bind request data to printer
	if err := c.Bind(&printer); err != nil {
		return c.JSON(http.StatusBadRequest, BaseResponse{
			Status:  false,
			Message: "Invalid request data",
			Data:    nil,
		})
	}

	// Check if printer with the same name already exists
	var existingPrinter Printer
	if err := DB.Where("nama = ?", printer.Name).First(&existingPrinter).Error; err == nil {
		return c.JSON(http.StatusConflict, BaseResponse{
			Status:  false,
			Message: "Printer with the same name already exists",
			Data:    nil,
		})
	}

	// Create new printer
	if result := DB.Create(&printer); result.Error != nil {
		return c.JSON(http.StatusInternalServerError, BaseResponse{
			Status:  false,
			Message: "Failed to create printer",
			Data:    nil,
		})
	}

	return c.JSON(http.StatusCreated, BaseResponse{
		Status:  true,
		Message: "Printer created successfully",
		Data:    printer,
	})
}

// GetPrinters retrieves a list of printers
func GetPrintersController(c echo.Context) error {
	var printers []Printer

	if result := DB.Order("name desc").Find(&printers); result.Error != nil {
		return c.JSON(http.StatusInternalServerError, BaseResponse{
			Status:  false,
			Message: "Failed to retrieve printers",
			Data:    nil,
		})
	}

	return c.JSON(http.StatusOK, BaseResponse{
		Status:  true,
		Message: "Printers retrieved successfully",
		Data:    printers,
	})
}

// UpdatePrinter updates an existing printer record
func UpdatePrinterController(c echo.Context) error {
	var printer Printer

	// Bind request data to printer
	if err := c.Bind(&printer); err != nil {
		return c.JSON(http.StatusBadRequest, BaseResponse{
			Status:  false,
			Message: "Invalid request data",
			Data:    nil,
		})
	}

	var existingPrinter Printer
	if err := DB.First(&existingPrinter, printer.ID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.JSON(http.StatusNotFound, BaseResponse{
				Status:  false,
				Message: "Printer not found",
				Data:    nil,
			})
		}
		return c.JSON(http.StatusInternalServerError, BaseResponse{
			Status:  false,
			Message: "Failed to find printer",
			Data:    nil,
		})
	}

	if result := DB.Model(&existingPrinter).Updates(printer).Error; result != nil {
		return c.JSON(http.StatusInternalServerError, BaseResponse{
			Status:  false,
			Message: "Failed to update printer",
			Data:    nil,
		})
	}

	return c.JSON(http.StatusOK, BaseResponse{
		Status:  true,
		Message: "Printer updated successfully",
		Data:    existingPrinter,
	})
}

// SoftDeletePrinter performs a soft delete on a printer record
func SoftDeletePrinterController(c echo.Context) error {
	id := c.Param("id")

	var printer Printer
	if err := DB.First(&printer, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.JSON(http.StatusNotFound, BaseResponse{
				Status:  false,
				Message: "Printer not found",
				Data:    nil,
			})
		}
		return c.JSON(http.StatusInternalServerError, BaseResponse{
			Status:  false,
			Message: "Failed to find printer",
			Data:    nil,
		})
	}

	if err := DB.Delete(&printer).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, BaseResponse{
			Status:  false,
			Message: "Failed to delete printer",
			Data:    nil,
		})
	}

	return c.JSON(http.StatusOK, BaseResponse{
		Status:  true,
		Message: "Printer soft-deleted successfully",
		Data:    nil,
	})
}

func RestorePrinterController(c echo.Context) error {
	id := c.Param("id")

	var printer Printer
	// Find the printer by ID including soft-deleted records
	if err := DB.Unscoped().Where("id = ?", id).First(&printer).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.JSON(http.StatusNotFound, BaseResponse{
				Status:  false,
				Message: "Printer not found",
				Data:    nil,
			})
		}
		return c.JSON(http.StatusInternalServerError, BaseResponse{
			Status:  false,
			Message: "Failed to find printer: " + err.Error(),
			Data:    nil,
		})
	}

	// Check if the printer is already active
	if printer.DeletedAt.Time.IsZero() {
		return c.JSON(http.StatusConflict, BaseResponse{
			Status:  false,
			Message: "Printer is not deleted",
			Data:    nil,
		})
	}

	// Restore the printer by setting DeletedAt to zero value
	if err := DB.Model(&printer).Update("DeletedAt", gorm.DeletedAt{}).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, BaseResponse{
			Status:  false,
			Message: "Failed to restore printer: " + err.Error(),
			Data:    nil,
		})
	}

	return c.JSON(http.StatusOK, BaseResponse{
		Status:  true,
		Message: "Printer restored successfully",
		Data:    printer,
	})
}

// DeletePrinter permanently deletes a printer record
func DeletePrinterController(c echo.Context) error {
	id := c.Param("id")

	var printer Printer
	// Find the printer by ID including soft-deleted records
	if err := DB.Unscoped().Where("id = ?", id).First(&printer).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.JSON(http.StatusNotFound, BaseResponse{
				Status:  false,
				Message: "Printer not found",
				Data:    nil,
			})
		}
		return c.JSON(http.StatusInternalServerError, BaseResponse{
			Status:  false,
			Message: "Failed to find printer: " + err.Error(),
			Data:    nil,
		})
	}

	// Perform hard delete
	if err := DB.Unscoped().Delete(&printer).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, BaseResponse{
			Status:  false,
			Message: "Failed to permanently delete printer: " + err.Error(),
			Data:    nil,
		})
	}

	return c.JSON(http.StatusOK, BaseResponse{
		Status:  true,
		Message: "Printer permanently deleted successfully",
		Data:    nil,
	})
}

// controller meja
func AddMejaController(c echo.Context) error {
	var meja Meja

	if err := c.Bind(&meja); err != nil {
		return c.JSON(http.StatusBadRequest, BaseResponse{
			Status:  false,
			Message: "Invalid request data",
			Data:    nil,
		})
	}

	if meja.Nama == "" {
		return c.JSON(http.StatusBadRequest, BaseResponse{
			Status:  false,
			Message: "Nama is required",
			Data:    nil,
		})
	}

	var existingMeja Meja
	if err := DB.Where("nama = ?", meja.Nama).First(&existingMeja).Error; err == nil {
		return c.JSON(http.StatusConflict, BaseResponse{
			Status:  false,
			Message: "Table with nama " + meja.Nama + " already exists",
			Data:    nil,
		})
	}

	if err := DB.Create(&meja).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, BaseResponse{
			Status:  false,
			Message: "Failed to create meja",
			Data:    nil,
		})
	}

	return c.JSON(http.StatusCreated, BaseResponse{
		Status:  true,
		Message: "Meja created successfully",
		Data:    meja,
	})
}
func GetMejasController(c echo.Context) error {
	var mejaList []Meja

	if err := DB.Find(&mejaList).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, BaseResponse{
			Status:  false,
			Message: "Failed to retrieve meja",
			Data:    nil,
		})
	}

	return c.JSON(http.StatusOK, BaseResponse{
		Status:  true,
		Message: "Meja retrieved successfully",
		Data:    mejaList,
	})
}
func UpdateMejaController(c echo.Context) error {
	id := c.Param("id")
	var updatedMeja Meja

	if err := c.Bind(&updatedMeja); err != nil {
		return c.JSON(http.StatusBadRequest, BaseResponse{
			Status:  false,
			Message: "Invalid request data",
			Data:    nil,
		})
	}

	if updatedMeja.Nama == "" {
		return c.JSON(http.StatusBadRequest, BaseResponse{
			Status:  false,
			Message: "Nama is required",
			Data:    nil,
		})
	}

	var existingMeja Meja
	if err := DB.First(&existingMeja, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.JSON(http.StatusNotFound, BaseResponse{
				Status:  false,
				Message: "Meja not found",
				Data:    nil,
			})
		}
		return c.JSON(http.StatusInternalServerError, BaseResponse{
			Status:  false,
			Message: "Failed to retrieve meja",
			Data:    nil,
		})
	}

	existingMeja.Nama = updatedMeja.Nama
	if err := DB.Save(&existingMeja).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, BaseResponse{
			Status:  false,
			Message: "Failed to update meja",
			Data:    nil,
		})
	}

	return c.JSON(http.StatusOK, BaseResponse{
		Status:  true,
		Message: "Meja updated successfully",
		Data:    existingMeja,
	})
}
func SoftDeleteMejaController(c echo.Context) error {
	id := c.Param("id")

	var meja Meja
	if err := DB.First(&meja, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.JSON(http.StatusNotFound, BaseResponse{
				Status:  false,
				Message: "Meja not found",
				Data:    nil,
			})
		}
		return c.JSON(http.StatusInternalServerError, BaseResponse{
			Status:  false,
			Message: "Failed to find meja",
			Data:    nil,
		})
	}

	if err := DB.Delete(&meja).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, BaseResponse{
			Status:  false,
			Message: "Failed to soft delete meja",
			Data:    nil,
		})
	}

	return c.JSON(http.StatusOK, BaseResponse{
		Status:  true,
		Message: "Meja soft-deleted successfully",
		Data:    nil,
	})
}
func RestoreMejaController(c echo.Context) error {
	id := c.Param("id")

	var meja Meja
	if err := DB.Unscoped().First(&meja, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.JSON(http.StatusNotFound, BaseResponse{
				Status:  false,
				Message: "Meja not found",
				Data:    nil,
			})
		}
		return c.JSON(http.StatusInternalServerError, BaseResponse{
			Status:  false,
			Message: "Failed to find meja",
			Data:    nil,
		})
	}

	if meja.DeletedAt.Time.IsZero() {
		return c.JSON(http.StatusConflict, BaseResponse{
			Status:  false,
			Message: "Meja is not deleted",
			Data:    nil,
		})
	}

	if err := DB.Model(&meja).Update("DeletedAt", nil).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, BaseResponse{
			Status:  false,
			Message: "Failed to restore meja",
			Data:    nil,
		})
	}

	return c.JSON(http.StatusOK, BaseResponse{
		Status:  true,
		Message: "Meja restored successfully",
		Data:    meja,
	})
}
func DeleteMejaController(c echo.Context) error {
	id := c.Param("id")

	if err := DB.Unscoped().Delete(&Meja{}, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.JSON(http.StatusNotFound, BaseResponse{
				Status:  false,
				Message: "Meja not found",
				Data:    nil,
			})
		}
		return c.JSON(http.StatusInternalServerError, BaseResponse{
			Status:  false,
			Message: "Failed to delete meja",
			Data:    nil,
		})
	}

	return c.JSON(http.StatusOK, BaseResponse{
		Status:  true,
		Message: "Meja permanently deleted successfully",
		Data:    nil,
	})
}

// add product menu
func CreateProductController(c echo.Context) error {
	var product Product

	// Bind request data to product
	if err := c.Bind(&product); err != nil {
		return c.JSON(http.StatusBadRequest, BaseResponse{
			Status:  false,
			Message: "Invalid request data",
			Data:    nil,
		})
	}

	// Check if product with the same name already exists
	//var existingProduct Product
	//if err := DB.Where("name = ?", product.Name).First(&existingProduct).Error; err == nil {
	//	formattedMessage := fmt.Sprintf("Product with name %s already exists", product.Name)
	//	return c.JSON(http.StatusConflict, BaseRespose{
	//		Status:  false,
	//		Message: formattedMessage,
	//		Data:    nil,
	//	})
	//}

	// Create new product
	result := DB.Create(&product)
	if result.Error != nil {
		return c.JSON(http.StatusInternalServerError, BaseResponse{
			Status:  false,
			Message: "Failed to add product",
			Data:    nil,
		})
	}

	return c.JSON(http.StatusCreated, BaseResponse{
		Status:  true,
		Message: "Successfully added product",
		Data:    product,
	})
}
func GetProductsController(c echo.Context) error {
	var products []Product

	result := DB.Find(&products)

	if result.Error != nil {
		return c.JSON(http.StatusInternalServerError, BaseResponse{
			Status:  false,
			Message: "Failed to retrieve products: " + result.Error.Error(),
			Data:    nil,
		})
	}

	return c.JSON(http.StatusOK, BaseResponse{
		Status:  true,
		Message: "Successfully retrieved products",
		Data:    products,
	})
}

// UpdateProductController updates an existing Product
func UpdateProductController(c echo.Context) error {
	id := c.Param("id")
	var updatedProduct Product

	// Bind request data to updatedProduct
	if err := c.Bind(&updatedProduct); err != nil {
		return c.JSON(http.StatusBadRequest, BaseResponse{
			Status:  false,
			Message: "Invalid request data",
			Data:    nil,
		})
	}

	// Validate updatedProduct data
	if updatedProduct.Name == "" {
		return c.JSON(http.StatusBadRequest, BaseResponse{
			Status:  false,
			Message: "Invalid product data: Name is required",
			Data:    nil,
		})
	}
	if updatedProduct.Varian == "" {
		return c.JSON(http.StatusBadRequest, BaseResponse{
			Status:  false,
			Message: "Invalid product data: Varian is required",
			Data:    nil,
		})
	}
	if updatedProduct.Price <= 0 {
		return c.JSON(http.StatusBadRequest, BaseResponse{
			Status:  false,
			Message: "Invalid product data: Price must be greater than zero",
			Data:    nil,
		})
	}

	// Find the existing product
	var existingProduct Product
	if err := DB.First(&existingProduct, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.JSON(http.StatusNotFound, BaseResponse{
				Status:  false,
				Message: "Product not found",
				Data:    nil,
			})
		}
		return c.JSON(http.StatusInternalServerError, BaseResponse{
			Status:  false,
			Message: "Failed to retrieve product",
			Data:    nil,
		})
	}

	// Update product details
	existingProduct.Category = updatedProduct.Category
	existingProduct.Name = updatedProduct.Name
	existingProduct.Varian = updatedProduct.Varian
	existingProduct.Price = updatedProduct.Price

	// Save the updated product
	if result := DB.Save(&existingProduct); result.Error != nil {
		return c.JSON(http.StatusInternalServerError, BaseResponse{
			Status:  false,
			Message: "Failed to update product",
			Data:    nil,
		})
	}

	return c.JSON(http.StatusOK, BaseResponse{
		Status:  true,
		Message: "Product updated successfully",
		Data:    existingProduct,
	})
}

// SoftDeleteProductController performs a soft delete of a Product
func SoftDeleteProductController(c echo.Context) error {
	id := c.Param("id")

	var product Product
	if err := DB.First(&product, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.JSON(http.StatusNotFound, BaseResponse{
				Status:  false,
				Message: "Product not found",
				Data:    nil,
			})
		}
		return c.JSON(http.StatusInternalServerError, BaseResponse{
			Status:  false,
			Message: "Failed to find product",
			Data:    nil,
		})
	}

	if err := DB.Delete(&product).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, BaseResponse{
			Status:  false,
			Message: "Failed to delete product",
			Data:    nil,
		})
	}

	return c.JSON(http.StatusOK, BaseResponse{
		Status:  true,
		Message: "Successfully soft-deleted product",
		Data:    nil,
	})
}

// RestoreProductController handles restoring a soft-deleted product
func RestoreProductController(c echo.Context) error {
	id := c.Param("id")

	var product Product
	if err := DB.Unscoped().First(&product, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.JSON(http.StatusNotFound, BaseResponse{
				Status:  false,
				Message: "Product not found",
				Data:    nil,
			})
		}
		return c.JSON(http.StatusInternalServerError, BaseResponse{
			Status:  false,
			Message: "Failed to find product",
			Data:    nil,
		})
	}

	if product.DeletedAt.Time.IsZero() {
		return c.JSON(http.StatusConflict, BaseResponse{
			Status:  false,
			Message: "Product is not deleted",
			Data:    nil,
		})
	}

	if err := DB.Model(&product).Update("DeletedAt", gorm.DeletedAt{}).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, BaseResponse{
			Status:  false,
			Message: "Failed to restore product",
			Data:    nil,
		})
	}

	return c.JSON(http.StatusOK, BaseResponse{
		Status:  true,
		Message: "Successfully restored product",
		Data:    product,
	})
}

// DeleteProductController handles hard deletion of a Product by ID
func DeleteProductController(c echo.Context) error {
	id := c.Param("id")

	var product Product
	// Find the product by ID
	if err := DB.Unscoped().First(&product, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.JSON(http.StatusNotFound, BaseResponse{
				Status:  false,
				Message: "Product not found",
				Data:    nil,
			})
		}
		return c.JSON(http.StatusInternalServerError, BaseResponse{
			Status:  false,
			Message: "Failed to find product: " + err.Error(),
			Data:    nil,
		})
	}

	// Perform hard delete
	if err := DB.Delete(&product).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, BaseResponse{
			Status:  false,
			Message: "Failed to delete product: " + err.Error(),
			Data:    nil,
		})
	}

	return c.JSON(http.StatusOK, BaseResponse{
		Status:  true,
		Message: "Product deleted successfully",
		Data:    nil,
	})
}

// CreateOrder handles creating an order and assigning it to printers
func CreateOrder(c echo.Context) error {
	var request CreateOrderRequest
	if err := c.Bind(&request); err != nil {
		return createErrorResponse(c, http.StatusBadRequest, "Invalid request payload")
	}

	tx := DB.Begin()
	if tx.Error != nil {
		return createErrorResponse(c, http.StatusInternalServerError, "Failed to start transaction")
	}
	defer func() {
		if r := recover(); r != nil || tx.Error != nil {
			tx.Rollback()
		}
	}()

	// Create the order
	order := Order{
		TableNumber: request.TableNumber,
		Status:      0,
	}
	if err := tx.Create(&order).Error; err != nil {
		tx.Rollback()
		return createErrorResponse(c, http.StatusInternalServerError, "Failed to create order")
	}

	// Process items and handle printers as before
	responsePrinters, debugInfo := processOrderItems(tx, request.Items, order.ID)

	if err := tx.Commit().Error; err != nil {
		return createErrorResponse(c, http.StatusInternalServerError, "Failed to commit transaction")
	}

	// Respond with the order details and printers
	return c.JSON(http.StatusCreated, BaseResponse{
		Status:  true,
		Message: "Order created successfully",
		Data: map[string]interface{}{
			"order_id":     order.ID,
			"table_number": order.TableNumber,
			"printers":     responsePrinters,
			"debug_info":   debugInfo,
		},
	})
}

func GetOrderController(c echo.Context) error {
	id := c.Param("id")
	var order Order
	if err := DB.Unscoped().First(&order, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return createErrorResponse(c, http.StatusNotFound, "Order not found")
		}
		return createErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve order")
	}

	return c.JSON(http.StatusOK, BaseResponse{
		Status:  true,
		Message: "Order retrieved successfully",
		Data:    order,
	})
}

func UpdateOrderController(c echo.Context) error {
	id := c.Param("id")

	var request UpdateOrderRequest
	if err := c.Bind(&request); err != nil {
		return createErrorResponse(c, http.StatusBadRequest, "Invalid request payload")
	}

	var order Order
	if err := DB.First(&order, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return createErrorResponse(c, http.StatusNotFound, "Order not found")
		}
		return createErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve order")
	}

	// Update order fields
	order.TableNumber = request.TableNumber
	order.Status = request.Status

	if err := DB.Save(&order).Error; err != nil {
		return createErrorResponse(c, http.StatusInternalServerError, "Failed to update order")
	}

	return c.JSON(http.StatusOK, BaseResponse{
		Status:  true,
		Message: "Order updated successfully",
		Data:    order,
	})
}

func SoftDeleteOrderController(c echo.Context) error {
	id := c.Param("id")
	var order Order
	if err := DB.First(&order, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return createErrorResponse(c, http.StatusNotFound, "Order not found")
		}
		return createErrorResponse(c, http.StatusInternalServerError, "Failed to find order")
	}

	if err := DB.Model(&order).Update("DeletedAt", gorm.DeletedAt{}).Error; err != nil {
		return createErrorResponse(c, http.StatusInternalServerError, "Failed to delete order")
	}

	return c.JSON(http.StatusOK, BaseResponse{
		Status:  true,
		Message: "Order soft-deleted successfully",
		Data:    nil,
	})
}

func RestoreOrderController(c echo.Context) error {
	id := c.Param("id")
	var order Order
	if err := DB.Unscoped().First(&order, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return createErrorResponse(c, http.StatusNotFound, "Order not found")
		}
		return createErrorResponse(c, http.StatusInternalServerError, "Failed to find order")
	}

	if order.DeletedAt.Time.IsZero() {
		return createErrorResponse(c, http.StatusConflict, "Order is not deleted")
	}

	if err := DB.Model(&order).Update("DeletedAt", nil).Error; err != nil {
		return createErrorResponse(c, http.StatusInternalServerError, "Failed to restore order")
	}

	return c.JSON(http.StatusOK, BaseResponse{
		Status:  true,
		Message: "Order restored successfully",
		Data:    order,
	})
}

func DeleteOrderController(c echo.Context) error {
	id := c.Param("id")
	if err := DB.Unscoped().Delete(&Order{}, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return createErrorResponse(c, http.StatusNotFound, "Order not found")
		}
		return createErrorResponse(c, http.StatusInternalServerError, "Failed to delete order")
	}

	return c.JSON(http.StatusOK, BaseResponse{
		Status:  true,
		Message: "Order deleted permanently",
		Data:    nil,
	})
}

// Helper function to handle error responses
func createErrorResponse(c echo.Context, statusCode int, message string) error {
	return c.JSON(statusCode, BaseResponse{
		Status:  false,
		Message: message,
		Data:    nil,
	})
}

// Function to process order items and handle printers
func processOrderItems(tx *gorm.DB, items []OrderItemRequest, orderID uint) (map[string][]string, []string) {
	// Map categories to printer names
	printerMap := map[string]string{
		"Minuman": "Printer Bar",
		"Makanan": "Printer Dapur",
	}

	// Fetch all printers once
	var allPrinters []Printer
	if err := tx.Find(&allPrinters).Error; err != nil {
		tx.Rollback()
		log.Printf("Failed to find printers: %v", err)
		return nil, []string{"Failed to find printers"}
	}

	// Create a map of printer names to IDs
	printerIDs := make(map[string][]string)
	for _, printer := range allPrinters {
		if printer.Name != "" { // Ensure the printer has a valid name
			printerIDs[printer.Name] = append(printerIDs[printer.Name], printer.ID)
		}
	}

	// Initialize response data and debug information
	responsePrinters := make(map[string][]string)
	var debugInfo []string

	for _, itemRequest := range items {
		var product Product
		if err := tx.Where("id = ?", itemRequest.ProductID).First(&product).Error; err != nil {
			tx.Rollback()
			return nil, []string{"Product not found"}
		}

		orderItem := OrderItem{
			OrderID:   orderID,
			ProductID: product.ID,
			Quantity:  itemRequest.Quantity,
		}
		if err := tx.Create(&orderItem).Error; err != nil {
			tx.Rollback()
			return nil, []string{"Failed to create order item"}
		}

		printerName, exists := printerMap[product.Category]
		if exists {
			ids, found := printerIDs[printerName]
			if !found {
				debugInfo = append(debugInfo, fmt.Sprintf("No printers found for category %s", printerName))
				log.Printf("No printers found for category %s", printerName)
				continue
			}
			log.Printf("Found printers for category %s: %+v", printerName, ids)
			if _, ok := responsePrinters[printerName]; !ok {
				responsePrinters[printerName] = ids
			}
			for _, printerID := range ids {
				orderPrinter := OrderPrinter{
					OrderID:   orderID,
					PrinterID: printerID,
				}
				if err := tx.Create(&orderPrinter).Error; err != nil {
					tx.Rollback()
					return nil, []string{"Failed to assign printer"}
				}
				log.Printf("Assigned printer %s to order %d", printerID, orderID)
			}
		} else {
			debugInfo = append(debugInfo, fmt.Sprintf("No printer mapping found for product category %s", product.Category))
			log.Printf("No printer mapping found for product category %s", product.Category)
		}
	}

	return responsePrinters, debugInfo
}

func GetMejaController(c echo.Context) error {
	id := c.Param("id")
	var meja Meja

	if err := DB.First(&meja, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return createErrorResponse(c, http.StatusNotFound, "Meja not found")
		}
		return createErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve meja")
	}

	return c.JSON(http.StatusOK, BaseResponse{
		Status:  true,
		Message: "Meja retrieved successfully",
		Data:    meja,
	})
}

// GetBill retrieves and calculates the total bill for a given table
func GetBill(c echo.Context) error {
	tableNumber := c.Param("table_number")

	var orders []Order
	if err := DB.Preload("Items.Product").Where("table_number = ?", tableNumber).Find(&orders).Error; err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to retrieve orders")
	}

	totalAmount := calculateTotalAmount(orders)

	response := struct {
		Orders      []Order `json:"orders"`
		TotalAmount float64 `json:"total_amount"`
	}{
		Orders:      orders,
		TotalAmount: totalAmount,
	}

	return c.JSON(http.StatusOK, response)
}

// calculateTotalAmount calculates the total amount of the given orders
func calculateTotalAmount(orders []Order) float64 {
	totalAmount := 0.0
	for _, order := range orders {
		for _, item := range order.Items {
			totalAmount += float64(item.Quantity) * item.Product.Price
		}
	}
	return totalAmount
}

// CalculateTotal calculates the total amount for a specific order ID including discounts
func CalculateTotal(orderID uint, db *gorm.DB) (float64, error) {
	var order Order
	if err := db.Preload("Items.Product").First(&order, orderID).Error; err != nil {
		return 0, err
	}

	total := calculateTotalAmount([]Order{order})

	var promos []Promo
	if err := db.Find(&promos).Error; err != nil {
		return 0, err
	}

	for _, promo := range promos {
		if containsAllProducts(order.Items, promo.ProductIDs) {
			total -= promo.Harga
		}
	}

	return total, nil
}

// containsAllProducts checks if all specified product IDs are present in the order items
func containsAllProducts(items []OrderItem, productIDs []uint) bool {
	productMap := make(map[uint]bool)
	for _, item := range items {
		productMap[item.ProductID] = true
	}

	for _, id := range productIDs {
		if !productMap[id] {
			return false
		}
	}

	return true
}

//func loadEnv() {
//	err := godotenv.Load()
//	if err != nil {
//		panic("Failed load env file")
//	}
//}
