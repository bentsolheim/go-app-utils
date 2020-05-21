package go_app_utils

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"net/http"
	"os"
	"regexp"
	"strings"
)

type Response struct {
	Message string
	Items   interface{}
}

func JsonResponse(w http.ResponseWriter, f func() (interface{}, error)) {
	encoder := json.NewEncoder(w)
	data, err := f()
	var response Response
	if err != nil {
		response = Response{Message: err.Error()}
	} else {
		response = Response{Message: "OK", Items: data}
	}
	if err := encoder.Encode(conventionalMarshaller{response}); err != nil {
		println(err.Error())
	}
}

var keyMatchRegex = regexp.MustCompile(`\"(\w+)\":`)

type conventionalMarshaller struct {
	Value interface{}
}

func (m conventionalMarshaller) MarshalJSON() ([]byte, error) {
	marshalled, err := json.Marshal(m.Value)

	converted := keyMatchRegex.ReplaceAllFunc(
		marshalled,
		func(match []byte) []byte {
			args := string(match)
			return []byte((args[:1] + strings.ToLower(string(args[1])) + args[2:]))
		},
	)

	return converted, err
}

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
		user:     getEnvOrDefault("DB_USER", defaults.user),
		password: getEnvOrDefault("DB_PASSWORD", defaults.password),
		host:     getEnvOrDefault("DB_HOST", defaults.host),
		port:     getEnvOrDefault("DB_PORT", defaults.port),
		name:     getEnvOrDefault("DB_NAME", defaults.name),
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

func getEnvOrDefault(name string, defaultValue string) string {
	value, isSet := os.LookupEnv(name)
	if !isSet {
		value = defaultValue
	}
	return value
}
