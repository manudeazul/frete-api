package repository

import (
	"database/sql"
	"fmt"
	"log"
	"superfrete-api/model"
)

type QuoteRepository struct {
	connection *sql.DB
}

func NewQuoteRepository(connection *sql.DB) QuoteRepository {
	return QuoteRepository{
		connection: connection,
	}
}

func (q *QuoteRepository) GetLastQuotes(limit int) ([]model.Quote, error) {
	var quotes []model.Quote
	query := "SELECT name, service, deadline, price FROM carrier ORDER BY id DESC"

	if limit != -1 {
		query = fmt.Sprintf("%s LIMIT %d", query, limit)
	}

	rows, err := q.connection.Query(query)
	if err != nil {
		log.Printf("Erro ao executar query: %v", err)
		return []model.Quote{}, err
	}
	defer rows.Close()

	for rows.Next() {
		var quote model.Quote
		if err := rows.Scan(&quote.Name, &quote.Service, &quote.Deadline, &quote.Price); err != nil {
			log.Printf("Erro ao escanear resultado: %v", err)
			continue
		}
		quotes = append(quotes, quote)
	}

	if err != nil {
		log.Printf("Erro nos resultados da query: %v", err)
		return []model.Quote{}, err
	}

	return quotes, nil
}

func (q *QuoteRepository) CreateQuote(quotes []model.Quote) error {
	for _, quote := range quotes {
		query := `INSERT INTO carrier (name, service, deadline, price) VALUES ($1, $2, $3, $4)`

		_, err := q.connection.Exec(query, quote.Name, quote.Service, quote.Deadline, quote.Price)
		if err != nil {
			log.Printf("Erro no insert: %v", err)
			return err
		}
	}

	return nil
}
