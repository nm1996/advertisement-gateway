package controllers

import (
	"gateway/interfaces"
	"gateway/model"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AdvertisementController struct {
	advertisementService interfaces.AdvertisementService
	logger               *log.Logger
}

func NewAdvertisementController(service interfaces.AdvertisementService, logger *log.Logger) interfaces.AdvertisementController {
	return &AdvertisementController{
		advertisementService: service,
		logger:               logger,
	}
}

func (controller *AdvertisementController) ReceiveAdvertisements(context *gin.Context) {
	var advertisements []model.Advertisement

	if err := context.BindJSON(&advertisements); err != nil {
		controller.logger.Printf("Failed to bind JSON: %v", err)
		context.AbortWithStatus(http.StatusBadRequest)
		return
	}

	go func() {
		controller.logger.Printf("Converting finished successfully! Processing to fill the queue!")
		controller.advertisementService.AddToQueue(&advertisements)
	}()

	context.Status(http.StatusOK)
}

func (controller *AdvertisementController) GetAdvertisements(context *gin.Context) {
	var data = controller.advertisementService.GetFromQueue()
	context.JSON(http.StatusOK, data)
}
