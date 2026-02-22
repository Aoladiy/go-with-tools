package config

import (
	"errors"
	"os"
	"strconv"
)

type Config struct {
	AppHost        string
	AppPort        int
	DbHost         string
	DbPort         string
	DbUsername     string
	DbPassword     string
	DbDatabaseName string
	DbSchema       string
	JwtSecret      string
}

func (c *Config) LoadEnv() error {
	appHost, appHostExists := os.LookupEnv("APP_HOST")
	appPort, appPortExists := os.LookupEnv("APP_PORT")
	dbHost, dbHostExists := os.LookupEnv("DB_HOST")
	dbPort, dbPortExists := os.LookupEnv("DB_PORT")
	dbUsername, dbUsernameExists := os.LookupEnv("DB_USERNAME")
	dbPassword, dbPasswordExists := os.LookupEnv("DB_PASSWORD")
	db, dbExists := os.LookupEnv("DB_DATABASE")
	dbSchema, dbSchemaExists := os.LookupEnv("DB_SCHEMA")
	jwtSecret, jwtSecretExists := os.LookupEnv("JWT_SECRET")

	if !appHostExists {
		return errors.New("APP_HOST .env isn't set")
	}
	if !appPortExists {
		return errors.New("APP_PORT .env isn't set")
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

	intAppPort, err := strconv.Atoi(appPort)
	if err != nil {
		return err
	}

	c.AppHost = appHost
	c.AppPort = intAppPort
	c.DbHost = dbHost
	c.DbPort = dbPort
	c.DbUsername = dbUsername
	c.DbPassword = dbPassword
	c.DbDatabaseName = db
	c.DbSchema = dbSchema
	c.JwtSecret = jwtSecret
	return nil
}
