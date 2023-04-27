package client

import (
	"fmt"

	raketapb "github.com/vanyaio/raketa-backend/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func NewGrpcClient(serverIP string, serverPort uint16) (raketapb.RaketaServiceClient, error) {
	conn, err := grpc.Dial(
		fmt.Sprintf("%s:%d", serverIP, serverPort),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, err
	}

	return raketapb.NewRaketaServiceClient(conn), nil
}
