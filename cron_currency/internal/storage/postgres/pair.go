package postgres

import (
	"context"

	"scrapper.go/cron_currency/internal/model"
	"scrapper.go/cron_currency/internal/storage"
	postgresql "scrapper.go/cron_currency/pkg/db"
)

type pairStorage struct {
	client postgresql.Client
}

func NewPairRepository(client postgresql.Client) storage.PairStorage {
	return &pairStorage{
		client: client,
	}
}

// AddPair implements storage.CurrencyStorage.
func (c *pairStorage) AddPair(ctx context.Context, base string, quote string) error {
	q := `INSERT INTO subscribed_pairs  (base_currency, quote_currency) 
	VALUES ($1, $2)`

	_, err := c.client.Exec(ctx, q, base, quote)
	if err != nil {
		return err
	}
	return nil

}

// GetAllPairs implements storage.CurrencyStorage.
func (c *pairStorage) GetAllPairs(ctx context.Context) ([]model.Pair, error) {
	q := `SELECT id, base_currency, quote_currency FROM subscribed_pairs`
	row, err := c.client.Query(ctx, q)
	if err != nil {
		return nil, err
	}

	pairs := make([]model.Pair, 0)
	for row.Next() {
		var prs model.Pair
		err := row.Scan(&prs.ID, &prs.Base, &prs.Quote)
		if err != nil {
			return nil, err
		}
		pairs = append(pairs, prs)
	}
	if err = row.Err(); err != nil {
		return nil, err
	}
	return pairs, nil

}

func (c *pairStorage) GetPairID(ctx context.Context, pair model.Pair) (int64, error) {
	var id int
	q := `SELECT id FROM subscribed_pairs WHERE base_currency = $1 AND quote_currency = $2 `

	qRow := c.client.QueryRow(ctx, q, pair.Base, pair.Quote)
	if err := qRow.Scan(&id); err != nil {
		return int64(id), err
	}
	return int64(id), nil

}
