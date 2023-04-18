package services

import (
	"encoding/json"
	"gateway/interfaces"
	"gateway/model"
	"gateway/utils"
	"log"

	"github.com/nats-io/nats.go"
)

type EventService struct {
	logger            *log.Logger
	Queue             *utils.SafeQueue
	natsConfiguration *utils.NatsConfiguration
}

func NewEventService(logger *log.Logger, queue *utils.SafeQueue, natsConfiguration *utils.NatsConfiguration) interfaces.EventService {
	return &EventService{
		logger:            logger,
		Queue:             queue,
		natsConfiguration: natsConfiguration,
	}
}

func (service *EventService) PublishMessage() {
	nc, err := nats.Connect(service.natsConfiguration.PORT)
	if err != nil {
		service.logger.Println("Error connecting to NATS server:", err)
		return
	}

	defer nc.Close()

	go func() {
		advertisement := service.Queue.Dequeue().(model.Advertisement)

		advJson, err := json.Marshal(advertisement)
		if err != nil {
			service.logger.Println("Error serializing message:", err)
			return
		}

		err = nc.Publish("adv_msg", advJson)
		if err != nil {
			service.logger.Println("Error publishing message:", err)
		} else {
			service.logger.Println("Message published successfully.")
		}
	}()
}
