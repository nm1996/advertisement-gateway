package server

import (
	"gateway/controllers"
	"gateway/interfaces"
	"gateway/services"
	"gateway/utils"
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

const PORT = ":8080"

type Server struct {
	router                  *gin.Engine
	advertisementController *controllers.AdvertisementController
}

func NewServer(logger *log.Logger) *Server {
	//initialize dependencies

	queue := utils.NewSafeQueue()
	advService := services.NewAdvertisementService(logger, queue)
	advController := controllers.NewAdvertisementController(advService, logger)

	natsConfiguration := utils.GetNatsConfiguration()
	eventService := services.NewEventService(logger, queue, natsConfiguration)

	go sendEvents(eventService)

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
	server.router.GET("/get-queue-size", server.advertisementController.GetAdvertisementsCount)
}

func (server *Server) Start() {
	server.router.Run(PORT)
}

func sendEvents(eventService interfaces.EventService) {
	ticker := time.Tick(20 * time.Second)

	for range ticker {
		err := eventService.PublishMessage()
		if err != nil {
			log.Fatal("Error when sending to nats: ", err)
		}
	}
}
