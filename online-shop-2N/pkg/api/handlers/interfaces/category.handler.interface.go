package interfaces

import "github.com/gin-gonic/gin"

type CategoryHandler interface {
	GetAllCategories(ctx *gin.Context)
	SaveCategory(ctx *gin.Context)
	SaveSubCategory(ctx *gin.Context)
}
