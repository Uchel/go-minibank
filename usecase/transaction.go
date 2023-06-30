package usecase

import (
	"github.com/Uchel/go-minibank/models/dto"
	"github.com/Uchel/go-minibank/repositories"
)

type TrxUc interface {
	Transfer(trx *dto.TrxReq) string
	Topup(trx *dto.TrxReq) string
	Histories(account string) any
}

type trxUc struct {
	trxRepo repositories.TrxRepo
}

func NewTrxUc(trxRepo repositories.TrxRepo) TrxUc {
	return &trxUc{
		trxRepo: trxRepo,
	}
}

func (u trxUc) Transfer(trx *dto.TrxReq) string {
	return u.trxRepo.TrxTransfer(trx)
}
func (u trxUc) Topup(trx *dto.TrxReq) string {
	return u.trxRepo.TrxTopUp(trx)
}
func (u trxUc) Histories(account string) any {
	return u.trxRepo.GetTrxHistoryByAccount(account)
}
