package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type DocumentRequest struct {
	FileName string `json:"file_name"`
	Author string `json:"author"`
}

func main() {

	router := gin.Default()
	
	// Health Check (GET)

	router.GET("/api/health", func(context *gin.Context) {
		context.JSON(http.StatusOK, gin.H {
			"message": "SecondBrain Gateway is running smoothly!",
		})
	})

	// Menerima data JSON (POST)

	router.POST("/api/upload", func(context *gin.Context) {
		var docReq DocumentRequest

		if err := context.ShouldBindJSON(&docReq); err != nil {
			context.JSON(http.StatusBadRequest, gin.H {
				"error": "Invalid JSON Format or wrong data type",
			})
			return
		}

		fmt.Printf("New Document!\nFile Name: %s\nWriter: %s\n", docReq.FileName, docReq.Author)

		context.JSON(http.StatusAccepted, gin.H{
			"message": fmt.Sprintf("The Document %s from %s is accepted", docReq.FileName, docReq.Author),
			"status":  "Processing",
		})
	})

	fmt.Println("🚀 Gateway starting on port :8080...")
	router.Run(":8080")
}