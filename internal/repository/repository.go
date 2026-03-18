package repository

import (
	"context"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresRepository struct {
	pool *pgxpool.Pool
}

func NewRepository(pool *pgxpool.Pool) *PostgresRepository {
	return &PostgresRepository{
		pool: pool,
	}
}

func (r *PostgresRepository) InsertResult(ctx context.Context, req *InsertResultReq) error {
	query := `
		INSERT INTO unpack_results (request_id, input_string, unpacked_result)
		VALUES ($1, $2, $3)
	`
	_, err := r.pool.Exec(ctx, query, req.RequestID, req.InputString, req.UnpackedResult)
	if err != nil {
		return err
	}
	return nil
}

func (r *PostgresRepository) SelectByID(ctx context.Context, id uuid.UUID) ([]Result, error) {
	query := `
		SELECT request_id, input_string, unpacked_result
		FROM unpack_results
		WHERE request_id = $1
	`
	rows, err := r.pool.Query(ctx, query, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var results []Result
	for rows.Next() {
		var res Result
		err := rows.Scan(&res.RequestID, &res.InputString, &res.UnpackedResult)
		if err != nil {
			return nil, err
		}
		results = append(results, res)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return results, nil
}