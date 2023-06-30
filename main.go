package main

import (
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"

	"github.com/Uchel/go-minibank/routes"
	"github.com/Uchel/go-minibank/utils"
)

func main() {
	r := gin.Default()
	db := utils.ConnectDB()

	routes.AccountApi(r, db)
	routes.TrxApi(r, db)

	r.Run(":" + utils.DotEnv("PORT"))
}
