package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/fideligo/secondbrain-gateway/proto"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type DocumentRequestJSON struct {
	FileName string `json:"file_name"`
	Author string `json:"author"`
}

func main() {
	// setup koneksi gRPC ke python
	// assumption: python ai runs on port 50051
	fmt.Println("Connecting to Python AI Engine on port :50051...")
	conn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect to gRPC server : %v", err)
	}
	defer conn.Close()
	
	// membuat "service" client dari koneksi di atas
	brainClient := proto.NewBrainServiceClient(conn)

	// setup gin router
	router := gin.Default()

	// Health Check (GET)
		router.GET("/api/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H {
			"message": "SecondBrain Gateway is running smoothly!",
		})
	})

	// Upload Endpoint (POST)
	router.POST("/api/upload", func(c *gin.Context) {
		var jsonReq DocumentRequestJSON

		// accept model from user
		if err := c.ShouldBindJSON(&jsonReq); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON format"})
			return
		}

		// bungkus ke dalam model gRPC (brain.pb.go)
		grpcReq := &proto.DocumentRequest{
			FileName: jsonReq.FileName,
			Author: jsonReq.Author,
			Content: []byte("Ini simulasi isi file PDF berupa bytes"), // later will be replaced with actual file
		}

		// use gRPC service to send to python (brain_grpc.pb.go)
		// give it a 5 seconds timeout, if python doesn't respond then cancel the task
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		grpcResponse, err := brainClient.ProcessDocument(ctx, grpcReq)
		if err != nil {
			// if python is off or error
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to connect to AI Engine",
				"details": err.Error(),
			})
			return
		}

		// return response from python to frontend
		c.JSON(http.StatusOK, gin.H{
			"message": "Sukses diproses oleh AI Engine!",
			"ai_response": gin.H{
				"success": grpcResponse.Success,
				"message": grpcResponse.Message,
				"document_id": grpcResponse.DocumentId,
			},
		})

	})


	fmt.Println("🚀 Gateway starting on port :8080...")
	router.Run(":8080")
}