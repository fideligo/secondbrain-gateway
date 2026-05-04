package handler

import (
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/fideligo/secondbrain-gateway/internal/client"
	"github.com/fideligo/secondbrain-gateway/internal/model"
	"github.com/gin-gonic/gin"

	"gorm.io/gorm"
)

// API Class
type APIHandler struct {
	brainClient *client.BrainClient
	db			*gorm.DB
}

// Constructor API
func NewAPIHandler(brainClient *client.BrainClient, db *gorm.DB) *APIHandler {
	return &APIHandler{
		brainClient: brainClient,
		db: db,
	}
}

// handle json
func (h *APIHandler) UploadDocument(c *gin.Context) {
	// 1. Catch uploaded file from json "document"
	file, err := c.FormFile("document")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No document file provided"})
		return
	}

	author := c.PostForm("author")
	if author == "" {
		author = "Anonymous"
	}

	// 2. Create filepath and save physical file (currently in backend)
	uploadDir := "./uploads"
	if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create directory"})
		return
	}

	// Define the filePath here so it can be used later for DB
	filePath := filepath.Join(uploadDir, file.Filename)

	if err := c.SaveUploadedFile(file, filePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save physical file"})
		return
	}

	// 3. Read bytes of the file, for gRPC
	fileContent, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read file"})
		return
	}
	defer fileContent.Close()

	fileBytes, err := io.ReadAll(fileContent)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read file bytes"})
		return
	}

	// 4. Send to AI
	response, err := h.brainClient.ProcessDocument(file.Filename, author, fileBytes)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "AI Engine failed to process",
			"details": err.Error(),
		})
		return
	}

	// 5. Save to PostgreSQL
	newDoc := model.Document{
		FileName:   file.Filename,
		Author:     author,
		FilePath:   filePath,
		Summary:    response.Message, 
		UploadedAt: time.Now(),
	}

	if result := h.db.Create(&newDoc); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to save to database",
			"details": result.Error.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":     "Success! Document processed and saved permanently",
		"document_id": newDoc.ID,
		"ai_response": response,
	})

}

func (h *APIHandler) Chat(c *gin.Context) {
	var req struct {
		Query string `json:"query"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body",
		})
		return
	}
	
	answer, err := h.brainClient.Chat(req.Query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"answer": answer,
	})
}