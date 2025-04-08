package repository

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

const (
	host     = "postgres"
	port     = 5432
	user     = "postgres"
	password = "1234"
	dbname   = "postgres"
)

func ConnectDB() (*sql.DB, error) {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+"password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatalf("Erro ao conectar com o banco: %v", err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatalf("Erro ao conectar com o banco: %v", err)
	}

	query := `
CREATE TABLE IF NOT EXISTS carrier (
	id SERIAL PRIMARY KEY,
	name VARCHAR(50) NOT NULL,
	service VARCHAR(50),
	deadline INTEGER NOT NULL,
	price NUMERIC(10, 2) NOT NULL
);`

	_, err = db.Exec(query)
	if err != nil {
		log.Fatalf("Erro ao criar tabela: %v", err)
	}

	return db, nil
}
