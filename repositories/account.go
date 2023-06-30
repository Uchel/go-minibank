package repositories

import (
	"database/sql"
	"errors"
	"log"

	"github.com/rs/xid"

	"github.com/Uchel/go-minibank/models"
	"github.com/Uchel/go-minibank/models/dto"
	"github.com/Uchel/go-minibank/utils"
)

type AccountRepo interface {
	CreateCustomer(req *dto.AccountReq) string
	CreateAccount(req *dto.AccountReq, customerId string) (any, error)
	GetDataAccountByEmail(email string) (dto.AccountData, error)
}

type accountRepo struct {
	db *sql.DB
}

func NewAccount(db *sql.DB) AccountRepo {
	repo := new(accountRepo)
	repo.db = db
	return repo
}

// =====================================================  Service Register ======================================================================
// create customer
func (r *accountRepo) CreateCustomer(req *dto.AccountReq) string {

	customer := models.Customer{}
	customer.ID = xid.New().String()
	customer.Name = req.Name
	customer.Email = req.Email
	customer.Phone = req.Phone
	customer.Address = req.Address
	customer.Password = req.Password
	query := "insert into customer (id,name,email,phone,address,password) values ($1,$2,$3,$4,$5,$6)"

	_, exeErr := r.db.Exec(query, customer.ID, customer.Name, customer.Email, customer.Phone, customer.Address, customer.Password)

	if exeErr != nil {
		log.Println(exeErr)
		return "failed to create customer"
	}

	return customer.ID
}

// Create Account
func (r *accountRepo) CreateAccount(req *dto.AccountReq, customerId string) (any, error) {
	account := models.Account{}
	account.ID = xid.New().String()
	account.CustomerId = customerId
	account.AccountNumber = utils.AccountNumberGenerate()
	account.Balance = req.Balance
	query := "insert into account (id,customer_id,account_number,balance) values($1,$2,$3,$4)"
	_, exeErr := r.db.Exec(query, account.ID, account.CustomerId, account.AccountNumber, account.Balance)

	if exeErr != nil {
		log.Println(exeErr)
		return nil, errors.New("failed to create account")
	}

	return "Success to create Account", nil
}

//============================================================================================================================================================

//=============================================================Service Login, Get Data By Email that Login ==========================================

func (r *accountRepo) GetDataAccountByEmail(email string) (dto.AccountData, error) {
	var dataAccount dto.AccountData

	query := "select c.name,c.email,c.password,c.address,c.phone,a.account_number,a.balance from account as a join customer as c  on a.customer_id = c.id where c.email = $1"

	row := r.db.QueryRow(query, email)

	if err := row.Scan(
		&dataAccount.Name,
		&dataAccount.Email,
		&dataAccount.Password,
		&dataAccount.Address,
		&dataAccount.Phone,
		&dataAccount.AccountNumber,
		&dataAccount.Balance,
	); err != nil {
		log.Println(err)
	}

	if dataAccount.Email == "" {
		return dataAccount, errors.New("account not found")
	}

	return dataAccount, nil
}
