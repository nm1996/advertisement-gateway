package interfaces

import "gateway/model"

type EventService interface {
	PublishMessage() error
	CreateChunks() [][]model.Advertisement
}
