package postgres

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"tg_bot/config"
	"time"
)

type PostgresDB struct {
	pool *pgxpool.Pool
}

func NewWithConfig(ctx context.Context, pgUser, pgPassword string, config *config.PostgresConfig) (*PostgresDB, error) {
	return NewWithConnString(ctx, createConnectionString(pgUser, pgPassword, config))
}

func NewWithConnString(ctx context.Context, connString string) (*PostgresDB, error) {
	cfg, err := pgxpool.ParseConfig(connString)
	if err != nil {
		return nil, err
	}

	connPool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		return nil, err
	}

	if err = connPool.Ping(ctx); err != nil {
		return nil, err
	}
	return &PostgresDB{pool: connPool}, nil
}

func createConnectionString(pgUser, pgPassword string, config *config.PostgresConfig) string {
	return fmt.Sprintf("host=%s port=%d dbname=%s user=%s password=%s sslmode=%s",
		config.Host, config.Port, config.Database, pgUser, pgPassword, config.SSLMode)
}

func (p *PostgresDB) Close() {
	p.pool.Close()
}

func (p *PostgresDB) SaveTranslation(ctx context.Context, userID int64, originalText, translatedText string) error {
	_, err := p.pool.Exec(ctx, "INSERT INTO translations(user_id, original_text, translated_text) VALUES ($1, $2, $3)",
		userID, originalText, translatedText)
	return err
}

func (p *PostgresDB) GetUserHistory(ctx context.Context, userID int64) (string, error) {
	rows, err := p.pool.Query(ctx, "SELECT original_text, translated_text, timestamp FROM translations WHERE user_id = $1 ORDER BY timestamp DESC", userID)
	if err != nil {
		return "", err
	}
	defer rows.Close()

	var history string
	for rows.Next() {
		var originalText, translatedText string
		var timestamp time.Time
		err = rows.Scan(&originalText, &translatedText, &timestamp)
		if err != nil {
			return "", err
		}
		history += fmt.Sprintf("Original: %s\nTranslated: %s\nTime: %s\n\n", originalText, translatedText, timestamp.Format("2006-01-02 15:04:05"))
	}

	if err = rows.Err(); err != nil {
		return "", err
	}

	if history == "" {
		history = "История переводов пуста."
	}

	return history, nil
}
