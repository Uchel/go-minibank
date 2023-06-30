package usecase

import (
	"errors"

	"github.com/Uchel/go-minibank/models/dto"
	"github.com/Uchel/go-minibank/repositories"
)

type AccountUc interface {
	Register(req *dto.AccountReq) (any, error)
	FindDataAccountByEmail(email string) (dto.AccountData, error)
	FindByEmail(email string) (dto.AccountData, error)
}

type accountUc struct {
	accountRepo repositories.AccountRepo
}

func NewAccountUc(accountRepo repositories.AccountRepo) AccountUc {
	return &accountUc{
		accountRepo: accountRepo,
	}
}

// Register Account===============================================================================================
func (u *accountUc) Register(req *dto.AccountReq) (any, error) {
	customerId := u.accountRepo.CreateCustomer(req)
	if customerId == "failed to create customer" {
		return nil, errors.New(customerId)
	}
	return u.accountRepo.CreateAccount(req, customerId)
}

// Login, Logout, dan bisa akses data setelah login =================================================================
func (u *accountUc) FindDataAccountByEmail(email string) (dto.AccountData, error) {
	return u.accountRepo.GetDataAccountByEmail(email)
}

func (u *accountUc) FindByEmail(email string) (dto.AccountData, error) {
	return u.accountRepo.GetDataAccountByEmail(email)
}
