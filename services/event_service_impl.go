package services

import (
	"encoding/json"
	"gateway/interfaces"
	"gateway/model"
	"gateway/utils"
	"log"
	"math"
	"os"
	"sync"

	"github.com/nats-io/nats.go"
)

// maximum message size is 1 MB
const maxMessageSize = 1000000

// maximum chunk size is 100 KB
const maxChunkSize = 100000

type EventService struct {
	logger         *log.Logger
	Queue          *utils.SafeQueue
	mutex          sync.Mutex
	natsConnection *nats.Conn
}

func NewEventService(logger *log.Logger, queue *utils.SafeQueue, natsConfiguration *utils.NatsConfiguration) interfaces.EventService {

	nc, err := nats.Connect(natsConfiguration.PORT)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	return &EventService{
		logger:         logger,
		Queue:          queue,
		mutex:          sync.Mutex{},
		natsConnection: nc,
	}
}

func (service *EventService) PublishMessage() error {
	chunks := service.CreateChunks()

	errChan := make(chan error, len(chunks))

	for i := range chunks {
		chunkBytes, err := json.Marshal(chunks[i])
		if err != nil {
			service.logger.Fatalf("Converting to json failed: %v", err)
			return err
		}

		go func(chunkBytes []byte) {

			service.mutex.Lock()
			defer service.mutex.Unlock()

			err = service.natsConnection.Publish("adv_msg", chunkBytes)
			if err != nil {
				errChan <- err
				return
			}
			service.logger.Println("Message sent succsessfully")
			errChan <- nil
		}(chunkBytes)
	}

	for i := 0; i < len(chunks); i++ {
		err := <-errChan
		if err != nil {
			service.logger.Fatalf("Error on chunk %d, message: %v", i, err)
			return err
		}
	}
	return nil

}

func (service *EventService) CreateChunks() [][]model.Advertisement {
	numChunks := int(math.Ceil(float64(service.Queue.Len()) / float64(maxChunkSize)))
	chunks := make([][]model.Advertisement, numChunks)

	service.logger.Printf("Created %d number of chunks", len(chunks))

	for i := 0; i < numChunks; i++ {
		startIndex := i * maxChunkSize
		endIndex := int(math.Min(float64((i+1)*maxChunkSize), float64(service.Queue.Len())))

		chunks[i] = make([]model.Advertisement, endIndex-startIndex)
		for j := 0; j < len(chunks[i]); j++ {
			chunks[i][j] = service.Queue.Dequeue().(model.Advertisement)
		}
	}

	return chunks
}
