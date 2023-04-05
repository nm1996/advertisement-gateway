package interfaces

import "gateway/model"

type AdvertisementService interface {
	AddToQueue(advertisements *[]model.Advertisement)
	GetFromQueue() *[]model.Advertisement
}
