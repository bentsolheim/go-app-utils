package db

import (
	"database/sql"
	"fmt"
	"github.com/bentsolheim/go-app-utils/utils"
	"github.com/pkg/errors"
)

type DbConfig struct {
	User     string
	Password string
	Host     string
	Port     string
	Name     string
}

func (c DbConfig) ConnectString(passwordOverride string) string {
	password := c.Password
	if passwordOverride != "" {
		password = passwordOverride
	}
	connectTemplate := "%s:%s@(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local&multiStatements=true"
	connectString := fmt.Sprintf(connectTemplate, c.User, password, c.Host, c.Port, c.Name)
	return connectString
}

func ReadDbConfig(defaults DbConfig) DbConfig {
	return DbConfig{
		User:     utils.GetEnvOrDefault("DB_USER", defaults.User),
		Password: utils.GetEnvOrDefault("DB_PASSWORD", defaults.Password),
		Host:     utils.GetEnvOrDefault("DB_HOST", defaults.Host),
		Port:     utils.GetEnvOrDefault("DB_PORT", defaults.Port),
		Name:     utils.GetEnvOrDefault("DB_NAME", defaults.Name),
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
