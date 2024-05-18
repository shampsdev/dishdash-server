package pg

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v4"
)

type Config struct {
	User     string
	Password string
	Host     string
	Port     uint16
	Database string
}

func (cfg *Config) ConnString() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.Database,
	)
}

func NewPostgresDB(ctx context.Context, cfg Config) (*pgx.Conn, error) {
	conn, err := pgx.Connect(ctx, cfg.ConnString())
	if err != nil {
		return nil, err
	}
	err = conn.Ping(ctx)
	return conn, err
}
