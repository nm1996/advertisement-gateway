package server

import (
	"gateway/controllers"
	"gateway/services"
	"gateway/utils"
	"log"

	"github.com/gin-gonic/gin"
)

const PORT = ":8080"

type Server struct {
	router                  *gin.Engine
	advertisementController *controllers.AdvertisementController
}

func NewServer(logger *log.Logger) *Server {
	//initialize dependencies

	// natsConfiguration := utils.GetNatsConfiguration()
	queue := utils.NewSafeQueue()
	advService := services.NewAdvertisementService(logger, queue)
	advController := controllers.NewAdvertisementController(advService, logger)
	// eventService := services.NewEventService(logger, queue, natsConfiguration)

	//initialize router
	router := gin.Default()

	return &Server{
		router:                  router,
		advertisementController: advController.(*controllers.AdvertisementController),
	}
}

func (server *Server) SetupRoutes() {
	server.router.POST("/advertisements", server.advertisementController.ReceiveAdvertisements)
	server.router.GET("/get-advertisements", server.advertisementController.GetAdvertisements)
}

func (server *Server) Start() {
	server.router.Run(PORT)
}
