package routes

import (
	"database/sql"

	"github.com/gin-gonic/gin"

	"github.com/Uchel/go-minibank/controllers"
	"github.com/Uchel/go-minibank/middleware"
	"github.com/Uchel/go-minibank/repositories"
	"github.com/Uchel/go-minibank/usecase"
)

func TrxApi(router *gin.Engine, db *sql.DB) {
	//batas waktu token

	trxRepo := repositories.NewTrx(db)
	trxUsecase := usecase.NewTrxUc(trxRepo)
	trxCtrl := controllers.NewTrxController(trxUsecase) //Secret key untuk kebutuhan jwt

	routerTrx := router.Group("/minibank/trx")
	routerTrx.Use(middleware.AuthMiddleware())
	routerTrx.POST("/transfer", trxCtrl.Transfer)
	routerTrx.POST("/topup", trxCtrl.TopUp)
	routerTrx.GET("/histories", trxCtrl.Histories)
}
