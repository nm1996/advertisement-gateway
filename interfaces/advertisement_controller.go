package interfaces

import "github.com/gin-gonic/gin"

type AdvertisementController interface {
	ReceiveAdvertisements(context *gin.Context)
	GetAdvertisements(context *gin.Context)
}
