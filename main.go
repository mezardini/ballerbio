package main


import (
	"log"
	_ "github.com/lib/pq"
	"ballerbio/routes"
	"github.com/gin-gonic/gin"
)



func main() {
	gin.SetMode(gin.ReleaseMode)
	// Start the server
	log.Printf("Starting Gin server on localhost:8081")
	routes.RegisterRoutes()
}