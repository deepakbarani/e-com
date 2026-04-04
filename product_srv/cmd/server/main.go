package main

import (
	"fmt"
	"log"
	"net"
	"product_srv/api/repository"
	"product_srv/api/service"
	"product_srv/internals/config"
	"product_srv/internals/connection"
	"product_srv/pkg/pb"

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

	auditRepo := repository.NewAuditRepository(db)
	authRepo := repository.NewProductRepository(db, auditRepo)
	authService := service.NewProductService(authRepo, auditRepo)

	pb.RegisterProductServiceServer(grpcServer, authService)

	// Start the server to respond to the requests
	fmt.Println("<------------- Product Server is running ------------->")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalln("Failed to serve:", err)
		return
	}
}
