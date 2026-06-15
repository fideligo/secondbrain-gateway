package main

import (
	"fmt"
	"log"
	"os"

	"github.com/fideligo/secondbrain-gateway/internal/client"
	"github.com/fideligo/secondbrain-gateway/internal/database"
	"github.com/fideligo/secondbrain-gateway/internal/handler"
	"github.com/fideligo/secondbrain-gateway/proto"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"github.com/gin-contrib/cors"
)

type DocumentRequestJSON struct {
	FileName string `json:"file_name"`
	Author string `json:"author"`
}

func main() {

	db := database.InitDB()

	// Get gRPC URL from environment variable, fallback to localhost if not set
	brainServiceURL := os.Getenv("BRAIN_SERVICE_URL")
	if brainServiceURL == "" {
		brainServiceURL = "localhost:50051"
	}

	fmt.Printf("⏳ Connecting to Python AI Engine at %s...\n", brainServiceURL)
	conn, err := grpc.NewClient(brainServiceURL, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect to gRPC server : %v", err)
	}
	defer conn.Close()

	grpcClient := proto.NewBrainServiceClient(conn) 
	
	brainClient := client.NewBrainClient(grpcClient) 
	
	apiHandler := handler.NewAPIHandler(brainClient, db)

	// Router
	router := gin.Default()
	
	router.Use(cors.New(cors.Config{
        AllowOrigins:     []string{"http://localhost:3000"},
        AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
        AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
        ExposeHeaders:    []string{"Content-Length"},
        AllowCredentials: true,
    }))
	
	// Daftarkan fungsi milik si Kasir ke jalur HTTP
	router.GET("/api/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "Gateway Running!"})
	})
	router.POST("/api/upload", apiHandler.UploadDocument)
	router.POST("/api/notes", apiHandler.UploadNote)
	router.POST("/api/chat", apiHandler.Chat)

	router.GET("/api/documents", apiHandler.GetDocuments)
	router.DELETE("/api/documents/:id", apiHandler.DeleteDocument)

	fmt.Println("🚀 Gateway starting on port :8080...")
	router.Run(":8080")
}