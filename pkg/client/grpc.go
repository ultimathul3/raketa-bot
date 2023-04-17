package client

import (
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type GrpcClient struct {
}

func NewGrpcClient(serverIP string, serverPort uint16) (*GrpcClient, error) {
	_, err := grpc.Dial(
		fmt.Sprintf("%s:%d", serverIP, serverPort),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, err
	}

	// TODO: return BotServiceClient
	return nil, nil
}
