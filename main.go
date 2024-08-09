package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"

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

// Order represents an order with its items
type Order struct {
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
	ID        string         `gorm:"size:1;primaryKey"`
	Name      string         `gorm:"size:50;uniqueIndex"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
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
	ID        uint           `gorm:"primaryKey"`
	Nama      string         `gorm:"size:50;uniqueIndex"` // Ensures unique Nama
	DeletedAt gorm.DeletedAt `gorm:"index"`               // Soft delete
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

	e := echo.New()

	//route api Promo
	e.POST("/api/v1/discount", AddPromoController)
	e.GET("/api/v1/discount", GetPromoByIDController)
	e.PUT("/api/v1/discount", UpdatePromoController)
	e.DELETE("/api/v1/discount/:id", DeletePromoController)
	//route api Meja
	e.POST("/api/v1/table", AddMejaController)
	e.GET("/api/v1/table", GetMejaController)
	e.PUT("/api/v1/table", UpdateMejaController)
	e.DELETE("/api/v1/table/:id", SoftDeleteMejaController)
	e.PUT("/api/v1/table/restore/:id", RestoreMejaController)
	e.DELETE("/api/v1/table/:id", DeleteMejaController)
	//route api Printer
	e.POST("/api/v1/printers", CreatePrinterController)
	e.GET("api/v1/printers", GetPrintersController)
	e.PUT("/api/v1/printers", UpdatePrinterController)
	e.DELETE("/api/v1/printers/:id", SoftDeletePrinterController)
	e.PUT("/api/v1/printers/restore/:id", RestorePrinterController)
	e.DELETE("/api/v1/printers/:id", DeletePrinterController)
	//post menu
	e.POST("/api/v1/product", CreateProductController)
	e.GET("api/v1/product", GetProductsController)
	//post order
	e.POST("/api/v1/neworder", CreateOrder)
	//route api Get bill
	e.GET("/api/v1/bill/:table_number", GetBill)
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
		&Promo{},
		&Printer{},
		&Meja{},
		&Product{},
		&Order{},
		&OrderItem{},
		&OrderItemRequest{},
		&OrderPrinter{},
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

	if result := DB.Order("nama desc").Find(&printers); result.Error != nil {
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
	if err := DB.First(&printer, id).Error; err != nil {
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

	if err := DB.Delete(&printer).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, BaseResponse{
			Status:  false,
			Message: "Failed to delete printer",
			Data:    nil,
		})
	}

	return c.JSON(http.StatusOK, BaseResponse{
		Status:  true,
		Message: "Printer deleted successfully",
		Data:    nil,
	})
}

// RestorePrinter restores a soft-deleted printer record
func RestorePrinterController(c echo.Context) error {
	id := c.Param("id")

	var printer Printer
	if err := DB.Unscoped().First(&printer, id).Error; err != nil {
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

	if printer.DeletedAt.Time.IsZero() {
		return c.JSON(http.StatusConflict, BaseResponse{
			Status:  false,
			Message: "Printer is not deleted",
			Data:    nil,
		})
	}

	if err := DB.Model(&printer).Update("DeletedAt", nil).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, BaseResponse{
			Status:  false,
			Message: "Failed to restore printer",
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
	if err := DB.First(&printer, id).Error; err != nil {
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

	if err := DB.Unscoped().Delete(&printer).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, BaseResponse{
			Status:  false,
			Message: "Failed to permanently delete printer",
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
// AddMejaController handles creating a new Meja entry
func AddMejaController(c echo.Context) error {
	var meja Meja

	// Bind request data to meja
	if err := c.Bind(&meja); err != nil {
		return c.JSON(http.StatusBadRequest, BaseResponse{
			Status:  false,
			Message: "Invalid request data",
			Data:    nil,
		})
	}

	// Validate meja data
	if meja.Nama == "" {
		return c.JSON(http.StatusBadRequest, BaseResponse{
			Status:  false,
			Message: "Invalid meja data: Nama is required",
			Data:    nil,
		})
	}

	// Check if promo with the same name already exists
	var existingMeja Meja
	if err := DB.Where("nama = ?", meja.Nama).First(&existingMeja).Error; err == nil {
		// If no error and a record is found
		formattedMessage := fmt.Sprintf("Table with nama %s already exists", meja.Nama)
		return c.JSON(http.StatusConflict, BaseResponse{
			Status:  false,
			Message: formattedMessage,
			Data:    nil,
		})
	}

	// Create new meja
	result := DB.Create(&meja)
	if result.Error != nil {
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

// GetMejaController retrieves a Meja by ID

// UpdateMejaController updates an existing Meja entry
func UpdateMejaController(c echo.Context) error {
	id := c.Param("id")
	var updatedMeja Meja

	// Bind request data to updatedMeja
	if err := c.Bind(&updatedMeja); err != nil {
		return c.JSON(http.StatusBadRequest, BaseResponse{
			Status:  false,
			Message: "Invalid request data",
			Data:    nil,
		})
	}

	// Validate updatedMeja data
	if updatedMeja.Nama == "" {
		return c.JSON(http.StatusBadRequest, BaseResponse{
			Status:  false,
			Message: "Invalid meja data: Nama is required",
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

	// Update meja details
	existingMeja.Nama = updatedMeja.Nama
	result := DB.Save(&existingMeja)
	if result.Error != nil {
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

// softdelete meja
func SoftDeleteMejaController(c echo.Context) error {
	// Extract category ID from request
	id := c.Param("id")

	// Find the category by ID
	var meja Meja
	if err := DB.First(&meja, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.JSON(http.StatusNotFound, BaseResponse{
				Status:  false,
				Message: "meja not found",
				Data:    nil,
			})
		}
		return c.JSON(http.StatusInternalServerError, BaseResponse{
			Status:  false,
			Message: "Failed to find category",
			Data:    nil,
		})
	}

	// Perform soft delete
	if err := DB.Delete(&meja).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, BaseResponse{
			Status:  false,
			Message: "Failed to delete meja",
			Data:    nil,
		})
	}

	return c.JSON(http.StatusOK, BaseResponse{
		Status:  true,
		Message: "Successfully deleted meja",
		Data:    nil,
	})
}

// restore meja
func RestoreMejaController(c echo.Context) error {
	// Extract category ID from request
	id := c.Param("id")

	// Find the category by ID (including soft-deleted records)
	var meja Meja
	if err := DB.Unscoped().First(&meja, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.JSON(http.StatusNotFound, BaseResponse{
				Status:  false,
				Message: "meja not found",
				Data:    nil,
			})
		}
		return c.JSON(http.StatusInternalServerError, BaseResponse{
			Status:  false,
			Message: "Failed to find meja",
			Data:    nil,
		})
	}

	// Check if the category is already active
	if meja.DeletedAt.Time.IsZero() {
		return c.JSON(http.StatusConflict, BaseResponse{
			Status:  false,
			Message: "meja is not deleted",
			Data:    nil,
		})
	}

	// Restore the category by setting DeletedAt to nil
	if err := DB.Model(&meja).Update("DeletedAt", nil).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, BaseResponse{
			Status:  false,
			Message: "Failed to restore meja",
			Data:    nil,
		})
	}

	return c.JSON(http.StatusOK, BaseResponse{
		Status:  true,
		Message: "Successfully restored meja",
		Data:    meja,
	})
}

// DeleteMejaController handles deleting a Meja by ID
func DeleteMejaController(c echo.Context) error {
	id := c.Param("id")

	// Soft delete the meja by ID
	result := DB.Delete(&Meja{}, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
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
		Message: "Meja deleted successfully",
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

// CreateOrder handles the creation of a new order
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
