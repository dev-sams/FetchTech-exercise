package main

import (
	"math"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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

var receipts = make(map[string]int)

func main() {
	router := gin.Default()

	router.POST("/receipts/process", processReceipt)

	router.GET("/receipts/:id/points", getPoints)

	router.Run() // Start the Gin server
}

func processReceipt(c *gin.Context) {
	var receipt Receipt
	if err := c.ShouldBindJSON(&receipt); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	receiptID := uuid.New().String()

	points := calculatePoints(receipt)

	receipts[receiptID] = points

	c.JSON(http.StatusOK, gin.H{"id": receiptID})
}

func getPoints(c *gin.Context) {
	id := c.Param("id")

	if points, exists := receipts[id]; exists {
		c.JSON(http.StatusOK, gin.H{"points": points})
	} else {
		c.JSON(http.StatusNotFound, gin.H{"error": "Receipt not found"})
	}
}

// Points calculation function implementing the rules
func calculatePoints(receipt Receipt) int {
	points := 0

	points += countAlphanumeric(receipt.Retailer)

	total, _ := strconv.ParseFloat(receipt.Total, 64)
	if total == math.Floor(total) {
		points += 50
	}

	if total*100 == math.Floor(total*100) && int(total*100)%25 == 0 {
		points += 25
	}

	points += (len(receipt.Items) / 2) * 5

	// Item-specific rules
	for _, item := range receipt.Items {
		if len(strings.TrimSpace(item.ShortDescription))%3 == 0 {
			price, _ := strconv.ParseFloat(item.Price, 64)
			points += int(math.Ceil(price * 0.2))
		}
	}

	date, _ := time.Parse("2006-01-02", receipt.PurchaseDate)
	if date.Day()%2 != 0 {
		points += 6
	}

	purchaseTime, _ := time.Parse("15:04", receipt.PurchaseTime)
	if purchaseTime.Hour() >= 14 && purchaseTime.Hour() < 16 {
		points += 10
	}

	return points
}

// Helper function to count alphanumeric characters
func countAlphanumeric(s string) int {
	alphanumeric := regexp.MustCompile(`[a-zA-Z0-9]`)
	return len(alphanumeric.FindAllString(s, -1))
}
