package api

import (
	"crypto-satangpro/middleware"
	"os"

	"github.com/gin-gonic/gin"
)


func InitGinFrameWork() {

	port := os.Getenv("PORT")

	//initz gin
	engine := gin.Default()
	engine.Use(middleware.CORSMiddleware())

	engine.POST("/transaction/list", GetTransactionListService)
	engine.Run(":" + port)

}
