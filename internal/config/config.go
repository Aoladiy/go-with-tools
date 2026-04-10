package config

import (
	"errors"
	"os"
	"strconv"
)

const (
	UserIdKey = "user_id"
)

type Config struct {
	AppHost string
	AppPort int

	LogLevel string

	AuthHost string
	AuthPort int

	DbHost         string
	DbPort         string
	DbUsername     string
	DbPassword     string
	DbDatabaseName string
	DbSchema       string

	JwtSecret string

	KafkaAddr string
}

func (c *Config) LoadEnv() error {
	appHost, appHostExists := os.LookupEnv("APP_HOST")
	appPort, appPortExists := os.LookupEnv("APP_PORT")
	authHost, authHostExists := os.LookupEnv("AUTH_HOST")
	authPort, authPortExists := os.LookupEnv("AUTH_PORT")

	logLevel, logLevelExists := os.LookupEnv("LOG_LEVEL")

	dbHost, dbHostExists := os.LookupEnv("DB_HOST")
	dbPort, dbPortExists := os.LookupEnv("DB_PORT")
	dbUsername, dbUsernameExists := os.LookupEnv("DB_USERNAME")
	dbPassword, dbPasswordExists := os.LookupEnv("DB_PASSWORD")
	db, dbExists := os.LookupEnv("DB_DATABASE")
	dbSchema, dbSchemaExists := os.LookupEnv("DB_SCHEMA")

	jwtSecret, jwtSecretExists := os.LookupEnv("JWT_SECRET")

	kafkaAddr, kafkaAddrExists := os.LookupEnv("KAFKA_ADDR")

	if !appHostExists {
		return errors.New("APP_HOST .env isn't set")
	}
	if !appPortExists {
		return errors.New("APP_PORT .env isn't set")
	}

	if !authHostExists {
		return errors.New("AUTH_HOST .env isn't set")
	}
	if !authPortExists {
		return errors.New("AUTH_PORT .env isn't set")
	}

	if !logLevelExists {
		return errors.New("LOG_LEVEL .env isn't set")
	}

	if !dbHostExists {
		return errors.New("DB_HOST .env isn't set")
	}
	if !dbPortExists {
		return errors.New("DB_PORT .env isn't set")
	}
	if !dbUsernameExists {
		return errors.New("DB_USERNAME .env isn't set")
	}
	if !dbPasswordExists {
		return errors.New("DB_PASSWORD .env isn't set")
	}
	if !dbExists {
		return errors.New("DB_DATABASE .env isn't set")
	}
	if !dbSchemaExists {
		return errors.New("DB_SCHEMA .env isn't set")
	}

	if !jwtSecretExists {
		return errors.New("JWT_SECRET .env isn't set")
	}

	if !kafkaAddrExists {
		return errors.New("KAFKA_ADDR .env isn't set")
	}

	intAppPort, err := strconv.Atoi(appPort)
	if err != nil {
		return err
	}

	intAuthPort, err := strconv.Atoi(authPort)
	if err != nil {
		return err
	}

	c.AppHost = appHost
	c.AppPort = intAppPort

	c.AuthHost = authHost
	c.AuthPort = intAuthPort

	c.LogLevel = logLevel

	c.DbHost = dbHost
	c.DbPort = dbPort
	c.DbUsername = dbUsername
	c.DbPassword = dbPassword
	c.DbDatabaseName = db
	c.DbSchema = dbSchema

	c.JwtSecret = jwtSecret

	c.KafkaAddr = kafkaAddr

	return nil
}
