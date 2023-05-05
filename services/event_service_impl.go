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
const maxMessageSize = 1_000_000

// maximum chunk size is 100 KB
const maxChunkSize = 1_000

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

	var currentMessageSize int
	var currentChunk []model.Advertisement

	for i := 0; i < service.Queue.Len(); i++ {
		item := service.Queue.Peek().(model.Advertisement)

		itemBytes, err := json.Marshal(item)
		if err != nil {
			service.logger.Printf("Error marshaling advertisement: %v", err)
			continue
		}

		if currentMessageSize+len(itemBytes) > maxMessageSize {
			// If adding the next item would exceed the max message size, start a new message.
			chunks = append(chunks, currentChunk)
			currentChunk = nil
			currentMessageSize = 0
		}

		if len(currentChunk) == maxChunkSize {
			// If the current chunk is at its max size, append it to the chunks slice and start a new chunk.
			chunks = append(chunks, currentChunk)
			currentChunk = nil
			currentMessageSize = 0
		}

		if currentChunk == nil {
			currentChunk = make([]model.Advertisement, 0)
		}

		currentChunk = append(currentChunk, item)
		service.Queue.Dequeue()

		currentMessageSize += len(itemBytes)
	}

	// Append the final chunk if there are any remaining items.
	if len(currentChunk) > 0 {
		chunks = append(chunks, currentChunk)
	}

	return chunks
}
