package api

import (
	"crypto-satangpro/db"
	"crypto-satangpro/middleware"
	"os"

	"github.com/gin-gonic/gin"
)


func InitGinFrameWork() {

	port := os.Getenv("PORT")

	//initz gin
	engine := gin.Default()
	engine.Use(middleware.CORSMiddleware())

	engine.GET("/healthz", db.Healthz)
	engine.POST("/transaction/list", GetTransactionListService)
	engine.POST("/address/add", AddAddressService)
	engine.Run(":" + port)

}
