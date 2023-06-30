package routes

import (
	"database/sql"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/Uchel/go-minibank/controllers"
	"github.com/Uchel/go-minibank/middleware"
	"github.com/Uchel/go-minibank/repositories"
	"github.com/Uchel/go-minibank/usecase"
	"github.com/Uchel/go-minibank/utils"
)

func AccountApi(router *gin.Engine, db *sql.DB) {
	//batas waktu token
	expireTime, _ := strconv.Atoi(utils.DotEnv("EXPIRE_TIME"))

	accountRepo := repositories.NewAccount(db)
	accountUsecase := usecase.NewAccountUc(accountRepo)
	accountCtrl := controllers.NewAccountController(accountUsecase, expireTime, utils.DotEnv("SECRET_KEY")) //Secret key untuk kebutuhan jwt

	router.POST("/auth/register", accountCtrl.Register)
	router.POST("/auth/login", accountCtrl.Login)
	routerAuth := router.Group("/minibank/auth")
	routerAuth.Use(middleware.AuthMiddleware())
	routerAuth.GET("/find-data", accountCtrl.FindByEmail)
	routerAuth.GET("/logout", accountCtrl.Logout)
}
