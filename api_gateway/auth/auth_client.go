package auth

import (
	"api-gateway/auth/pb"
	"api-gateway/internal/config"
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func NewClientService() pb.AuthServiceClient {

	//Getting the auth service URL from the enviroment variable
	authServiceURL := config.GetAuthserviceURL()

	//Creating the gRPC client for the auth service
	grpcClient, err := grpc.NewClient(authServiceURL, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		fmt.Println("Failed to get the auth *grpc.ClientConn : ", err)
	}

	return pb.NewAuthServiceClient(grpcClient)
}
