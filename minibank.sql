-- postgresql


-- enum
create type trx_type as enum(
	'topup',
	'transfer'
);

create table account(
id varchar(100) primary key,
customer_id varchar(100) not null ,
account_number varchar(100) not null,
balance integer not null
);

create table customer(
id varchar(100) primary key,
name varchar(100),
	email varchar(100),
	phone varchar(20),
	address varchar(200),
	password varchar(200) not null
);

create table transaction(
		id varchar(200) primary key,
		sender_account varchar(200),
		receiver_account varchar(200),
		trx_type trx_type,
		amount integer,
		created_at timestamp
);

create table transaction_report(
trx_id varchar(200),
my_account varchar(200),
kredit integer,
	debet integer,
	trx_type trx_type,
	account_x varchar(200),
	balance integer,
	created_at timestamp
);

ALTER table account ADD CONSTRAINT account_customer foreign key (customer_id) references customer(id);
alter table account alter column balance set default 200000;

alter table transaction_report alter column debet set default 0;
alter table transaction_report alter column kredit set default 0;

--untuk get detail tabel account dan customer

--select c.name,c.email,a.account_number,c.password,c.address,c.phone,a.balance 
--from account as a join customer as c  on a.customer_id = c.id;