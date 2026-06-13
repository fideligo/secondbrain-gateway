package client

import (
	"context"
	"time"

	"github.com/fideligo/secondbrain-gateway/proto"
)

type BrainClient struct {
	grpcClient proto.BrainServiceClient
}

type MessageHistory struct {
	Role string
	Content string
}

// NewBrainClient (constructor)
func NewBrainClient(grpcClient proto.BrainServiceClient) *BrainClient {
	return &BrainClient{
		grpcClient: grpcClient,
	}
}

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

func (c *BrainClient) Chat(query string, history []MessageHistory) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*300)
	defer cancel()

	var protoHistory []*proto.ChatMessage
	for _, msg := range history {
		protoHistory = append(protoHistory, &proto.ChatMessage{
			Role: msg.Role,
			Content: msg.Content,
		})
	}

	req := &proto.ChatRequest{
		Query: query,
		History: protoHistory,
	}

	res, err := c.grpcClient.Chat(ctx, req)
	if err != nil {
		return "", err
	}

	return res.Answer, nil
}

func (b *BrainClient) ProcessNote(title, author, content string) (*proto.DocumentResponse, error) {
	req := &proto.NoteRequest{
		Title:   title,
		Author:  author,
		Content: content,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	return b.grpcClient.ProcessNote(ctx, req)
}

