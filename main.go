package main

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/derickit/go-rest-api/internal/db"
	"github.com/derickit/go-rest-api/internal/logger"
	"github.com/derickit/go-rest-api/internal/models"
	"github.com/derickit/go-rest-api/internal/server"
	"github.com/derickit/go-rest-api/internal/util"
)

const (
	serviceName = "ecommerce-orders"
	defaultPort = "8080"
)

var version string

func main() {
	fmt.Println("Hello, World!")

}

func run() error {
	upTime := time.Now().UTC().Format(time.RFC3339)
	sigHandler := util.NewSignalHandler()

	svcEnv := MustEnvConfig()

	lgr := logger.Setup(svcEnv)
	dbCredentials, err := db.MongoDBCredentialFromSideCar(svcEnv.MongoVaultSideCar)
	if err != nil {
		lgr.Fatal().Err(err).Msg("failed to fetch db credentials")
		return err
	}
	connOpts := &db.ConnectionOpts{
		Database:     svcEnv.DBName,
		PrintQueries: svcEnv.PrintQueries,
	}

	dbConnMgr, err := db.NewMongoManager(dbCredentials, connOpts, lgr)
	if err != nil {
		lgr.Fatal().Err(err).Msg("unable to initialize db connection")
		return err
	}
	sigHandler.OnSignal(func() {
		dErr := dbConnMgr.Disconnect()
		if dErr != nil {
			lgr.Error().Err(dErr).Msg("unable to disconnect from db ,potential connection leak")
			return
		}
	})

	lgr.Info().Str("name", serviceName).Str("environment", svcEnv.Name).
		Str("started", upTime).Str("version", version).Msg("service details starting the service")

	server.StartService(svcEnv, dbConnMgr, lgr)
	lgr.Fatal().Msg("service stopped")
	return nil
}

func MustEnvConfig() models.ServiceEnv {
	envName := os.Getenv("environment")
	if envName == "" {
		envName = "local"
	}

	port := os.Getenv("port")
	if port == "" {
		port = defaultPort
	}

	dbName := os.Getenv("dbName")
	if dbName == "" {
		panic("dbName  should be defined in env configuration ")

	}

	printDBQueries, err := strconv.ParseBool(os.Getenv("printDBQueries"))
	if err != nil {
		printDBQueries = false
	}

	mongoSideCar := os.Getenv("mongoSideCar")
	if mongoSideCar == "" {
		panic("mongo sidecar file path should be defined in env configuration")
	}

	disableAuth, authEnvErr := strconv.ParseBool(os.Getenv("disableAuth"))
	if authEnvErr != nil {
		disableAuth = false
	}

	logLevel := os.Getenv("logLevel")
	if logLevel == "" {
		logLevel = "info"
	}

	envConfigurations := models.ServiceEnv{
		Name:              envName,
		Port:              port,
		PrintQueries:      printDBQueries,
		MongoVaultSideCar: mongoSideCar,
		DisableAuth:       disableAuth,
		DBName:            dbName,
		LogLevel:          logLevel,
	}
	return envConfigurations
}
