package controllers

import (
	"log"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"

	"github.com/Uchel/go-minibank/models/dto"
	"github.com/Uchel/go-minibank/usecase"
)

type AccountController struct {
	accountUc usecase.AccountUc
	waktu     int
	secret    string
}

func NewAccountController(accountUC usecase.AccountUc, waktu int, secret string) *AccountController {
	controller := AccountController{
		accountUc: accountUC,
		waktu:     waktu,
		secret:    secret,
	}
	return &controller
}

// ======================================AUTH=========================================================
// Register
func (c *AccountController) Register(ctx *gin.Context) {

	var newAccount dto.AccountReq

	if err := ctx.ShouldBind(&newAccount); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Input Field"})
		return
	}
	hashPass, err := bcrypt.GenerateFromPassword([]byte(newAccount.Password), 10)
	if err != nil {
		log.Println("Failed to hash password:", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	newAccount.Password = string(hashPass)

	res, err := c.accountUc.Register(&newAccount)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"response": res,
	})

}

// Login=======================
func (c *AccountController) Login(ctx *gin.Context) {
	var loginReq dto.LoginReq
	if err := ctx.ShouldBind(&loginReq); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Input Field"})
		return
	}

	dataAccount, err := c.accountUc.FindDataAccountByEmail(loginReq.Email)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid email or password"})
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(dataAccount.Password), []byte(loginReq.Password))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid email or password"})
		return
	}

	expire := time.Now().Add(time.Minute * time.Duration(c.waktu))
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["email"] = loginReq.Email
	claims["accNum"] = dataAccount.AccountNumber
	claims["exp"] = expire.Unix()

	tokenString, err := token.SignedString([]byte(c.secret))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "gagal generate token"})
		return
	}

	ctx.SetSameSite(http.SameSiteLaxMode)
	ctx.SetCookie("Authorization", tokenString, expire.Minute(), "/", "localhost", false, true)
	ctx.JSON(http.StatusOK, gin.H{"message": "Login Successfully"})
}

// Logout=====================================
func (c *AccountController) Logout(ctx *gin.Context) {
	expire := time.Now().Add(time.Minute * time.Duration(c.waktu))
	ctx.SetCookie("Authorization", "", -expire.Minute(), "/", "", false, true)

	ctx.JSON(http.StatusOK, gin.H{"message": "Logout Succesfully"})
}

//======================================================================================================================

// Service =============================================================Get Data By Email that Login ==========================================
func (c *AccountController) FindByEmail(ctx *gin.Context) {
	claims := ctx.MustGet("claims").(jwt.MapClaims)
	email := claims["email"].(string)
	accNum := claims["accNum"].(string)

	dataAccount, err := c.accountUc.FindDataAccountByEmail(email)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "email not found",
		})
		return
	}
	dataAccount.Password = ""

	ctx.JSON(http.StatusOK, gin.H{
		"data":   dataAccount,
		"accNum": accNum,
	})

}
