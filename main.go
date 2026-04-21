package main

import (
	"fmt"
	"log"

	"github.com/fideligo/secondbrain-gateway/internal/client"
	"github.com/fideligo/secondbrain-gateway/internal/handler"
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

	fmt.Println("⏳ Connecting to Python AI Engine on port :50051...")
	conn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect to gRPC server : %v", err)
	}
	defer conn.Close()

	grpcClient := proto.NewBrainServiceClient(conn) 
	
	brainClient := client.NewBrainClient(grpcClient) 
	
	apiHandler := handler.NewAPIHandler(brainClient) 

	// 3. BUKA TOKO (Router)
	router := gin.Default()
	
	// Daftarkan fungsi milik si Kasir ke jalur HTTP
	router.GET("/api/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "Gateway Running!"})
	})
	router.POST("/api/upload", apiHandler.UploadDocument)

	fmt.Println("🚀 Gateway starting on port :8080...")
	router.Run(":8080")
}