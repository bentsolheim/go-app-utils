package db

import (
	"database/sql"
	"fmt"
	"github.com/bentsolheim/go-app-utils/utils"
	"github.com/pkg/errors"
)

type DbConfig struct {
	user     string
	password string
	host     string
	port     string
	name     string
}

func (c DbConfig) ConnectString(passwordOverride string) string {
	password := c.password
	if passwordOverride != "" {
		password = passwordOverride
	}
	connectTemplate := "%s:%s@(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local"
	connectString := fmt.Sprintf(connectTemplate, c.user, password, c.host, c.port, c.name)
	return connectString
}

func ReadDbConfig(defaults DbConfig) DbConfig {
	return DbConfig{
		user:     utils.GetEnvOrDefault("DB_USER", defaults.user),
		password: utils.GetEnvOrDefault("DB_PASSWORD", defaults.password),
		host:     utils.GetEnvOrDefault("DB_HOST", defaults.host),
		port:     utils.GetEnvOrDefault("DB_PORT", defaults.port),
		name:     utils.GetEnvOrDefault("DB_NAME", defaults.name),
	}
}

func ConnectToDb(c DbConfig) (*sql.DB, error) {

	db, err := sql.Open("mysql", c.ConnectString(""))
	if err != nil {
		return nil, fmt.Errorf("unable to connect to database (%s): %v", c.ConnectString("***"), err)
	}
	if err := db.Ping(); err != nil {
		return nil, errors.Wrap(err, "error while verifying database connection")
	}

	return db, nil
}
