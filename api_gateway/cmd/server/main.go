package main

import (
	"api-gateway/internal/config"
	"api-gateway/routes"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {

	// loading the enviroment varibles
	err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
		return
	}

	//Initializing the gin
	g := gin.Default()

	// Added Recovery for the panic situtions
	g.Use(gin.Recovery())

	// Connected the routes
	routes.Routes(g)

	//Starting the Main server to listen request form the client
	fmt.Println("<------------- Gateway Server is running ------------->")
	if err := g.Run(config.GetPort()); err != nil {
		fmt.Println("Server error:", err)
		return
	}

}
