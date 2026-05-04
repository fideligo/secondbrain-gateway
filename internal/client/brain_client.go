package client

import (
	"context"
	"time"

	"github.com/fideligo/secondbrain-gateway/proto"
)

// pembungkus untuk layanan gRPC
type BrainClient struct {
	grpcClient proto.BrainServiceClient
}

// NewBrainClient (constructor)
func NewBrainClient(grpcClient proto.BrainServiceClient) *BrainClient {
	return &BrainClient{
		grpcClient: grpcClient,
	}
}

// ProcessDocument
func (b *BrainClient) ProcessDocument(fileName, author string, content []byte) (*proto.DocumentResponse, error) {

	req := &proto.DocumentRequest{
		FileName: fileName,
		Author: author,
		Content: content,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Second)
	defer cancel()

	return b.grpcClient.ProcessDocument(ctx, req)
}

func (c *BrainClient) Chat(query string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*300) // AI butuh waktu berpikir
	defer cancel()

	req := &proto.ChatRequest{
		Query: query,
	}

	res, err := c.grpcClient.Chat(ctx, req)
	if err != nil {
		return "", err
	}

	return res.Answer, nil
}