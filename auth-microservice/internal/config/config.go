package config

import (
	"errors"
	"os"
	"strconv"
	"time"
)

const (
	UserIdKey = "user_id"
)

type Config struct {
	AppPort int

	DbHost         string
	DbPort         string
	DbUsername     string
	DbPassword     string
	DbDatabaseName string
	DbSchema       string

	RdbAddr         string
	RdbUsr          string
	RdbPsw          string
	RdbId           int
	RdbMaxRetries   int
	RdbReadTimeout  time.Duration
	RdbWriteTimeout time.Duration
	RdbMinIdleConns int

	JwtSecret string
}

func (c *Config) LoadEnv() error {
	appPort, appPortExists := os.LookupEnv("APP_PORT")

	dbHost, dbHostExists := os.LookupEnv("DB_HOST")
	dbPort, dbPortExists := os.LookupEnv("DB_PORT")
	dbUsername, dbUsernameExists := os.LookupEnv("DB_USERNAME")
	dbPassword, dbPasswordExists := os.LookupEnv("DB_PASSWORD")
	db, dbExists := os.LookupEnv("DB_DATABASE")
	dbSchema, dbSchemaExists := os.LookupEnv("DB_SCHEMA")

	rdbAddr, rdbAddrExists := os.LookupEnv("REDIS_ADDR")
	rdbUsr, rdbUsrExists := os.LookupEnv("REDIS_USERNAME")
	rdbPsw, rdbPswExists := os.LookupEnv("REDIS_PASSWORD")
	rdbId, rdbIdExists := os.LookupEnv("REDIS_DB_IDENTIFIER")
	rdbMaxRetries, rdbMaxRetriesExists := os.LookupEnv("REDIS_MAX_RETRIES")
	rdbReadTimeout, rdbReadTimeoutExists := os.LookupEnv("REDIS_READ_TIMEOUT")
	rdbWriteTimeout, rdbWriteTimeoutExists := os.LookupEnv("REDIS_WRITE_TIMEOUT")
	rdbMinIdleConns, rdbMinIdleConnsExists := os.LookupEnv("REDIS_MIN_IDLE_CONNS")
	jwtSecret, jwtSecretExists := os.LookupEnv("JWT_SECRET")

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

	if !rdbAddrExists {
		return errors.New("REDIS_ADDR .env isn't set")
	}
	if !rdbUsrExists {
		return errors.New("REDIS_USERNAME .env isn't set")
	}
	if !rdbPswExists {
		return errors.New("REDIS_PASSWORD .env isn't set")
	}
	if !rdbIdExists {
		return errors.New("REDIS_DB_IDENTIFIER .env isn't set")
	}
	if !rdbMaxRetriesExists {
		return errors.New("REDIS_MAX_RETRIES .env isn't set")
	}
	if !rdbReadTimeoutExists {
		return errors.New("REDIS_READ_TIMEOUT .env isn't set")
	}
	if !rdbWriteTimeoutExists {
		return errors.New("REDIS_WRITE_TIMEOUT .env isn't set")
	}
	if !rdbMinIdleConnsExists {
		return errors.New("REDIS_MIN_IDLE_CONNS .env isn't set")
	}

	if !jwtSecretExists {
		return errors.New("JWT_SECRET .env isn't set")
	}

	intAppPort, err := strconv.Atoi(appPort)
	if err != nil {
		return err
	}

	intRdbId, err := strconv.Atoi(rdbId)
	if err != nil {
		return err
	}

	intRdbMaxRetries, err := strconv.Atoi(rdbMaxRetries)
	if err != nil {
		return err
	}

	intRdbReadTimeout, err := strconv.Atoi(rdbReadTimeout)
	if err != nil {
		return err
	}

	intRdbWriteTimeout, err := strconv.Atoi(rdbWriteTimeout)
	if err != nil {
		return err
	}

	intRdbMinIdleConns, err := strconv.Atoi(rdbMinIdleConns)
	if err != nil {
		return err
	}

	c.AppPort = intAppPort

	c.DbHost = dbHost
	c.DbPort = dbPort
	c.DbUsername = dbUsername
	c.DbPassword = dbPassword
	c.DbDatabaseName = db
	c.DbSchema = dbSchema

	c.RdbAddr = rdbAddr
	c.RdbUsr = rdbUsr
	c.RdbPsw = rdbPsw
	c.RdbId = intRdbId
	c.RdbMaxRetries = intRdbMaxRetries
	c.RdbReadTimeout = time.Duration(intRdbReadTimeout) * time.Second
	c.RdbWriteTimeout = time.Duration(intRdbWriteTimeout) * time.Second
	c.RdbMinIdleConns = intRdbMinIdleConns

	c.JwtSecret = jwtSecret

	return nil
}
