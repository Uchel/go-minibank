package controllers

import (
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"

	"github.com/Uchel/go-minibank/models/dto"
	"github.com/Uchel/go-minibank/usecase"
)

type TrxController struct {
	trxUc usecase.TrxUc
}

func NewTrxController(trxUc usecase.TrxUc) *TrxController {
	controller := TrxController{
		trxUc: trxUc,
	}
	return &controller
}

func (c TrxController) Transfer(ctx *gin.Context) {
	claims := ctx.MustGet("claims").(jwt.MapClaims)
	accNum := claims["accNum"].(string)
	var newTrx *dto.TrxReq

	if err := ctx.ShouldBind(&newTrx); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Input Field"})
		return
	}

	newTrx.SenderAccount = accNum
	newTrx.Type = "transfer"

	respon := c.trxUc.Transfer(newTrx)

	if respon == "internal server eror" {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": respon,
		})
		return
	}
	if respon == "account number not found" {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": respon,
		})

		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": respon,
	})

}
func (c TrxController) TopUp(ctx *gin.Context) {
	claims := ctx.MustGet("claims").(jwt.MapClaims)
	accNum := claims["accNum"].(string)
	var newTrx *dto.TrxReq

	if err := ctx.ShouldBind(&newTrx); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Input Field"})
		return
	}

	newTrx.SenderAccount = accNum
	newTrx.ReceiverAccount = accNum
	newTrx.Type = "topup"

	respon := c.trxUc.Topup(newTrx)

	if respon == "internal server eror" {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": respon,
		})
		return
	}
	if respon == "account number not found" {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": respon,
		})

		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": respon,
	})

}
func (c TrxController) Histories(ctx *gin.Context) {
	claims := ctx.MustGet("claims").(jwt.MapClaims)
	accNum := claims["accNum"].(string)

	respon := c.trxUc.Histories(accNum)

	if respon == "internal server eror" {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": respon,
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"message": respon,
	})

}
