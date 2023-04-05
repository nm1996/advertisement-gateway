package services

import (
	"gateway/interfaces"
	"gateway/model"
	"gateway/utils"
	"log"
)

type AdvertisementService struct {
	logger *log.Logger
	Queue  *utils.SafeQueue
}

func NewAdvertisementService(logger *log.Logger, queue *utils.SafeQueue) interfaces.AdvertisementService {
	return &AdvertisementService{
		logger: logger,
		Queue:  queue,
	}
}

func (service *AdvertisementService) AddToQueue(advertisements *[]model.Advertisement) {
	for _, advertise := range *advertisements {
		service.logger.Printf("Adding advertisement: %s to queue", advertise.Name)
		service.Queue.Enqueue(advertise)
		service.logger.Printf("Adding advertisement: %s to queue finished", advertise.Name)
	}
	service.logger.Printf("Adding all advertisements finished")
}

func (service *AdvertisementService) GetFromQueue() *[]model.Advertisement {
	var length = service.Queue.Len()

	result := make([]model.Advertisement, length)
	for i := 0; i < length; i++ {
		result[i] = service.Queue.Dequeue().(model.Advertisement)
	}

	return &result
}
