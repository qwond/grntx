package repository

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/qwond/grntx/internal/domain"
)

// mapRowToPair maps pgx.Row to Pair struct
func mapRowToPair(r pgx.Row) (domain.Pair, error) {
	var p domain.Pair
	err := r.Scan(
		&p.Pair,
		&p.AskUnit,
		&p.BidUnit,
		&p.MinAsk,
		&p.MinBid,
		&p.MakerFee,
		&p.TakerFee,
		&p.PricePrecision,
		&p.VolumePrecision,
		&p.CreatedAt,
		&p.UpdatedAt,
	)

	return p, err
}

const PairUpsertSQL = `
	INSERT INTO pairs (
		pair,
		ask_unit,
		bid_unit,
		min_ask,
		min_bid,
		maker_fee,
		taker_fee,
		price_precision,
		volume_precision
	) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	ON CONFLICT (pair)
	DO UPDATE SET
		ask_unit = $2,
		bid_unit = $3,
		min_ask = $4,
		min_bid = $5,
		maker_fee = $6,
		taker_fee = $7,
		price_precision = $8,
		volume_precision = $9,
		updated_at = EXTRACT(EPOCH FROM NOW())::BIGINT
	RETURNING (xmax = 0) AS is_insert`

// UpsertPair - update or insert pair entity returning bool,error
// where bool tells "it's new record" and error - db errors.
func (repo *Repository) PairUpsert(ctx context.Context, pr domain.Pair) (bool, error) {
	var isInsert bool
	err := repo.pool.QueryRow(ctx, PairUpsertSQL,
		pr.Pair,
		pr.AskUnit,
		pr.BidUnit,
		pr.MinAsk,
		pr.MinBid,
		pr.MakerFee,
		pr.TakerFee,
		pr.PricePrecision,
		pr.VolumePrecision,
	).Scan(&isInsert)
	if err != nil {
		return isInsert, err
	}

	return isInsert, nil
}

const PairsListSQL = `
	SELECT
		pair,
		ask_unit,
		bid_unit,
		min_ask,
		min_bid,
		maker_fee,
		taker_fee,
		price_precision,
		volume_precision,
		created_at,
		updated_at
	FROM pairs`

// PairsList returns all pairs stored in database.
func (repo *Repository) PairsList(ctx context.Context) ([]domain.Pair, error) {
	rows, err := repo.pool.Query(ctx, PairsListSQL)
	if err != nil {
		return nil, err
	}

	var pairs []domain.Pair
	for rows.Next() {
		pair, err := mapRowToPair(rows)
		if err != nil {
			return nil, err
		}
		pairs = append(pairs, pair)
	}

	return pairs, nil
}
