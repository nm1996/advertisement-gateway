package main

import (
	"gateway/server"
	"log"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	// initialize logger
	gin.SetMode(gin.DebugMode)
	logger := log.New(os.Stdout, "[myapp] ", log.LstdFlags|log.Lmicroseconds|log.LUTC)

	// create new server instance
	server := server.NewServer(logger)

	// set up routes
	server.SetupRoutes()

	// start server
	server.Start()
}
