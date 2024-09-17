package main

import (
    "github.com/gin-gonic/gin"
    "github.com/google/uuid"
    "net/http"
)

type Receipt struct {
    Retailer     string `json:"retailer" binding:"required"`
    PurchaseDate string `json:"purchaseDate" binding:"required"`
    PurchaseTime string `json:"purchaseTime" binding:"required"`
    Items        []Item `json:"items" binding:"required"`
    Total        string `json:"total" binding:"required"`
}

type Item struct {
    ShortDescription string `json:"shortDescription" binding:"required"`
    Price            string `json:"price" binding:"required"`
}

var receipts = make(map[string]int) // Store receipt points in-memory

func main() {
    router := gin.Default()

    // Register routes
    router.POST("/receipts/process", processReceipt)
    router.GET("/receipts/:id/points", getPoints)

    // Start the Gin server
    router.Run()
}

// processReceipt handles receipt processing and points calculation
func processReceipt(c *gin.Context) {
    var receipt Receipt

    // Bind JSON to Receipt struct and handle errors
    if err := c.ShouldBindJSON(&receipt); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    // Generate a unique ID for the receipt
    receiptID := uuid.New().String()

    // Calculate points based on receipt data
    points := calculatePoints(receipt)

    // Store receipt points in memory
    receipts[receiptID] = points

    // Return the receipt ID in the response
    c.JSON(http.StatusOK, gin.H{"id": receiptID})
}

// getPoints retrieves points for a given receipt ID
func getPoints(c *gin.Context) {
    receiptID := c.Param("id")

    if points, exists := receipts[receiptID]; exists {
        c.JSON(http.StatusOK, gin.H{"points": points})
    } else {
        c.JSON(http.StatusNotFound, gin.H{"error": "Receipt not found"})
    }
}

// calculatePoints calculates points for a given receipt
// This function should be implemented with the actual logic
func calculatePoints(receipt Receipt) int {
    // Implement points calculation logic based on receipt details
    return 0 // Placeholder
}