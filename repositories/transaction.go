package repositories

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/rs/xid"

	"github.com/Uchel/go-minibank/models/dto"
)

type TrxRepo interface {
	TrxTransfer(trx *dto.TrxReq) string
	TrxTopUp(trx *dto.TrxReq) string
	GetTrxHistoryByAccount(account string) any
}

type trxRepo struct {
	db *sql.DB
}

func NewTrx(db *sql.DB) TrxRepo {
	repo := new(trxRepo)
	repo.db = db
	return repo
}

// =======================================Validate untuk debugging transaction =============================
func Validate(err error, message string, tx *sql.Tx) {
	if err != nil {
		tx.Rollback()
		fmt.Println(err, "Transaction Rollback", message)
	}
}

// =========================================Transfer==============================================================
func (r *trxRepo) TrxTransfer(trx *dto.TrxReq) string {
	tx, err := r.db.Begin()
	if err != nil {
		log.Println(err)
		return "internal server eror"
	}
	trx.ID = xid.New().String()
	r.InsertTrx(trx, tx)
	//mengambil data pengirim dan penerima
	senderData := r.GetByAccountNumber(trx.SenderAccount, tx)
	receiverData := r.GetByAccountNumber(trx.ReceiverAccount, tx)
	balanceSender := senderData.Balance - trx.Amount
	balanceReceiver := receiverData.Balance + trx.Amount
	//Mengurangi balance pengirim dan Menambah balance Penerima sesuai amount
	r.UpdateBalanceByAccountNumber(balanceSender, senderData.AccountNumber, tx)     //update account sender : mengurangi balance
	r.UpdateBalanceByAccountNumber(balanceReceiver, receiverData.AccountNumber, tx) //update account receiver : menambah balance
	r.InsertTrxReportSender(trx, senderData.AccountNumber, receiverData.AccountNumber, balanceSender, tx)
	r.InsertTrxReportReceiver(trx, receiverData.AccountNumber, senderData.AccountNumber, balanceReceiver, tx)

	err = tx.Commit()
	if err != nil {
		log.Println(err)
		return "account number not found"
	} else {
		return "transaction successully"
	}
}

// =========================================Riwayat Transaksi==============================================================
func (r *trxRepo) GetTrxHistoryByAccount(account string) any {
	var histories []dto.TrxReport

	query := "select trx_id,kredit,debet,trx_type,account_x,balance,created_at from transaction_report where my_account = $1"
	rows, err := r.db.Query(query, account)
	if err != nil {
		log.Println(err)
	}
	defer rows.Close()

	for rows.Next() {
		var history dto.TrxReport

		if err := rows.Scan(&history.TrxId, &history.Kredit, &history.Debet, &history.Type, &history.AccountNumberXId, &history.Balance, &history.CreatedAt); err != nil {
			log.Println(err)
			return "internal server error"
		}
		histories = append(histories, history)
	}

	if len(histories) == 0 {
		return "no data"
	}
	return histories
}

// =========================================Top Up==============================================================
func (r *trxRepo) TrxTopUp(trx *dto.TrxReq) string {
	tx, err := r.db.Begin()
	if err != nil {
		log.Println(err)
		return "internal server eror"
	}

	trx.ID = xid.New().String()

	r.InsertTrx(trx, tx)

	//mengambil data account yang ingin di top Up
	senderData := r.GetByAccountNumber(trx.SenderAccount, tx)
	balanceSender := senderData.Balance + trx.Amount

	r.UpdateBalanceByAccountNumber(balanceSender, senderData.AccountNumber, tx) // menambah balance
	r.InsertTrxReportTopup(trx, senderData.AccountNumber, senderData.AccountNumber, balanceSender, tx)
	err = tx.Commit()
	if err != nil {
		log.Println(err)
		return "account number not found"
	} else {
		return "Top Up successully"
	}
}

// ==================================================================Transfer ==============================================================================================
func (r *trxRepo) InsertTrx(trx *dto.TrxReq, tx *sql.Tx) {
	query := "insert into transaction (id,sender_account,receiver_account,trx_type,amount,created_at) values($1,$2,$3,$4,$5,$6)"

	_, err := tx.Exec(query, trx.ID, trx.SenderAccount, trx.ReceiverAccount, trx.Type, trx.Amount, time.Now())
	Validate(err, "InsertTrx", tx)

}
func (r *trxRepo) UpdateBalanceByAccountNumber(balance int, accountNumber string, tx *sql.Tx) {
	query := "UPDATE account SET balance = $1 WHERE account_number = $2"
	_, err := r.db.Exec(query, balance, accountNumber)
	Validate(err, "UpdateBalanceByAccount", tx)
}

// Get data account sender atau receiver untuk proses transfer
func (r *trxRepo) GetByAccountNumber(accountNumber string, tx *sql.Tx) dto.AccountData {
	var dataAccount dto.AccountData

	query := "select c.name,c.email,c.password,c.address,c.phone,a.account_number,a.balance from account as a join customer as c  on a.customer_id = c.id where a.account_number = $1"

	row := tx.QueryRow(query, accountNumber)

	err := row.Scan(
		&dataAccount.Name,
		&dataAccount.Email,
		&dataAccount.Password,
		&dataAccount.Address,
		&dataAccount.Phone,
		&dataAccount.AccountNumber,
		&dataAccount.Balance,
	)

	Validate(err, "GetByAccountNumber", tx)

	return dataAccount
}

func (r *trxRepo) InsertTrxReportSender(trx *dto.TrxReq, acc1 string, acc2 string, balance int, tx *sql.Tx) {
	query := "insert into transaction_report (trx_id,my_account,account_x,debet,trx_type,balance,created_at) values($1,$2,$3,$4,$5,$6,$7)"
	_, err := tx.Exec(query, trx.ID, acc1, acc2, trx.Amount, trx.Type, balance, time.Now())
	Validate(err, "InsertTrxReportSender", tx)

}
func (r *trxRepo) InsertTrxReportReceiver(trx *dto.TrxReq, acc1 string, acc2 string, balance int, tx *sql.Tx) {
	query := "insert into transaction_report (trx_id,my_account,account_x,kredit,trx_type,balance,created_at) values($1,$2,$3,$4,$5,$6,$7)"
	_, err := tx.Exec(query, trx.ID, acc1, acc2, trx.Amount, trx.Type, balance, time.Now())
	Validate(err, "InsertTrxReportReceiver", tx)
}

// Get data account sender atau receiver untuk proses topup
func (r *trxRepo) InsertTrxReportTopup(trx *dto.TrxReq, acc1 string, acc2 string, balance int, tx *sql.Tx) {
	query := "insert into transaction_report (trx_id,my_account,account_x,kredit,trx_type,balance,created_at) values($1,$2,$3,$4,$5,$6,$7)"

	_, err := tx.Exec(query, trx.ID, acc1, acc2, trx.Amount, trx.Type, balance, time.Now())
	Validate(err, "InsertTrxReportSender", tx)

}
