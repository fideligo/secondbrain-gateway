package handler

import (
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"
	"fmt"

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
		History []client.MessageHistory `json:"history"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body",
		})
		return
	}
	
	answer, err := h.brainClient.Chat(req.Query, req.History)
	
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"answer": answer,
	})
}

func (h *APIHandler) GetDocuments(c *gin.Context) {
	var docs []model.Document

	if result := h.db.Order("uploaded_at desc").Find(&docs); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch documents",
			"details": result.Error.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Success retrieving documents",
		"data": docs,
	})
}

func (h *APIHandler) DeleteDocument(c *gin.Context) {
    id := c.Param("id")
    var doc model.Document

    if result := h.db.First(&doc, id); result.Error != nil {
        c.JSON(http.StatusNotFound, gin.H{
            "error": "Document not found",
        })
        return
    }

    err := os.Remove(doc.FilePath)
    if err != nil {
        fmt.Printf("Warning: Failed to delete physical file %s: %v\n", doc.FilePath, err)
    }

    if result := h.db.Delete(&doc); result.Error != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "error": "Failed to delete from database",
        })
        return 
    }

    c.JSON(http.StatusOK, gin.H{
        "message": fmt.Sprintf("Document with ID %s successfully deleted", id),
    })
}

