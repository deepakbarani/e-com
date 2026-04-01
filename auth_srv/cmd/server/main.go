package main

import (
	"auth_srv/api/repository"
	"auth_srv/api/service"
	"auth_srv/internals/config"
	"auth_srv/internals/connection"
	"auth_srv/pkg/pb"
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"
)

func main() {

	// Load the enviroment variables in the auth_srv
	err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	// Getting the database in the auth_srv
	db, err := connection.InitDB()
	if err != nil {
		log.Fatal(err)
		return
	}

	// Initialize server and  Listening
	lis, err := net.Listen("tcp", config.GetPort())
	if err != nil {
		log.Fatalln("failed at listening : ", err)
		return
	}

	//Create the gRPC server
	grpcServer := grpc.NewServer()

	authRepo := repository.NewAuthRepository(db)
	authService := service.NewAuthService(authRepo)

	pb.RegisterAuthServiceServer(grpcServer, authService)

	// Start the server to respond to the requests
	fmt.Println("<------------- AUth Server is running ------------->")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalln("Failed to serve:", err)
		return
	}
}
