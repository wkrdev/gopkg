package dbms

type DBCON struct {
	HOST string
	USER string
	PASS string
}

type RDBCON struct {
	DB   int
	HOST string
	PASS string
}

type Config struct {
	RWDB  *DBCON
	RODB  *DBCON
	REDIS *RDBCON
}
