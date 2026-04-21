package handler

import (
	"net/http"

	"github.com/fideligo/secondbrain-gateway/internal/client"
	"github.com/fideligo/secondbrain-gateway/internal/model"
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
	var jsonReq model.DocumentRequestJSON

	// accepts json
	if err := c.ShouldBindJSON(&jsonReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON format"})
		return
	}

	dummyContent := []byte("Ini simulasi isi file PDF berupa bytes")
	
	resp, err := h.brainClient.ProcessDocument(jsonReq.FileName, jsonReq.Author, dummyContent)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to connect to AI Engine",
			"details": err.Error(),
		})
		return
	}

	// Kasih struk balasan ke pelanggan
	c.JSON(http.StatusOK, gin.H{
		"message": "Sukses diproses oleh AI Engine!",
		"ai_response": gin.H{
			"success":     resp.Success,
			"message":     resp.Message,
			"document_id": resp.DocumentId,
		},
	})
}