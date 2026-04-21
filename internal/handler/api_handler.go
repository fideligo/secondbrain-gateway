package handler

import (
	"io"
	"net/http"

	"github.com/fideligo/secondbrain-gateway/internal/client"
	"github.com/gin-gonic/gin"
)

// API Class
type APIHandler struct {
	brainClient *client.BrainClient
}

// Constructor API
func NewAPIHandler(brainClient *client.BrainClient) *APIHandler {
	return &APIHandler{
		brainClient: brainClient,
	}
}

// handle json
func (h *APIHandler) UploadDocument(c *gin.Context) {

	// catch uploaded file named "document" from form-data
	file, err := c.FormFile("document")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No document file provided in the request"})
		return
	}

	// catch author name
	author := c.PostForm("author")
	if author == "" {
		author = "Anonymous"
	}

	// open file
	fileContent, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read file bytes"})
		return
	}

	// read file into raw bytes 
	fileBytes, err := io.ReadAll(fileContent)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read file bytes"})
		return
	}

	// send file bytes to the AI engine via gRPC
	response, err := h.brainClient.ProcessDocument(file.Filename, author, fileBytes)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "AI Engine failed to process the document",
			"details": err.Error(),
		})
		return
	}

	// Return success to the user
	c.JSON(http.StatusOK, gin.H{
		"message": "Success! Document processed by AI Engine",
		"ai_response": response,
	})
}