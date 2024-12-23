package client

import (
	"context"

	pb "github.com/minhthong582000/soa-404/api/v1/pb/random"
)

// Client is a simple client for the Random service.
type Client struct {
	randClient pb.RandomServiceClient
}

// NewClient creates a new client.
func NewClient(randClient pb.RandomServiceClient) *Client {
	return &Client{
		randClient: randClient,
	}
}

// GetRandNumber gets a random number from the server.
func (c Client) GetRandNumber(ctx context.Context, seed int64) (int64, error) {
	reply, err := c.randClient.GetRandNumber(ctx, &pb.GetRandNumberRequest{
		SeedNum: seed,
	})
	if err != nil {
		return -1, err
	}

	return reply.Number, nil
}
