package pg

import (
	"fmt"
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
