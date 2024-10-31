package main

import (
	"testing"

	"github.com/derickit/go-rest-api/internal/models"
	"github.com/stretchr/testify/assert"
)

func TestMustEnvConfig(t *testing.T) {
	t.Run("MissingEnvVariables", func(t *testing.T) {
		resetEnv(t)
		assert.Panics(t, func() {
			_ = MustEnvConfig()

		}, "MustEnvConfig did not panic with missing environment variables")
	})
	t.Run("ValidEnvVariables", func(t *testing.T) {
		resetEnv(t)
		t.Setenv("environment", "test")
		t.Setenv("port", "8080")
		t.Setenv("dbName", "testDB")
		t.Setenv("MongoVaultSideCar", "/path/to/mongo/sidecar")
		t.Setenv("logLevel", "debug")

		expectedConfig := models.ServiceEnv{
			Name:              "test",
			Port:              "8080",
			PrintQueries:      false,
			MongoVaultSideCar: "/path/to/mongo/sidecar",
			DisableAuth:       false,
			DBName:            "testDB",
			LogLevel:          "debug",
		}
		actualConfig := MustEnvConfig()
		assert.Equal(t, expectedConfig, actualConfig, "MustEnvConfig did not return the expected configuration")
	})
}

func TestMustEnvConfig_Default(t *testing.T) {
	t.Run("Default Env Values", func(t *testing.T) {
		resetEnv(t)

		expectedConfig := models.ServiceEnv{
			Name:              "test",
			Port:              defaultPort,
			PrintQueries:      false,
			MongoVaultSideCar: "/path/to/mongo/sidecar",
			DisableAuth:       false,
			DBName:            "testDB",
			LogLevel:          "info",
		}
		actualConfig := MustEnvConfig()
		assert.Equal(t, expectedConfig, actualConfig, "MustEnvConfig did not return the expected configuration")
	})
}

func TestMustEnvConfig_FailOnSideCar(t *testing.T) {
	t.Run("Fail on MongoSide Car", func(t *testing.T) {
		resetEnv(t)
		t.Setenv("environment", "test")
		t.Setenv("dbName", "testDB")
		assert.Panics(t, func() {
			_ = MustEnvConfig()

		}, "MustEnvConfig did not panic with missing mongo sidecar")
	})
}

func resetEnv(t *testing.T) {
	t.Setenv("environment", "")
	t.Setenv("port", "")
	t.Setenv("logLevel", "")
	t.Setenv("dbName", "")
	t.Setenv("MongoVaultSideCar", "")
	t.Setenv("printDBQueries", "")

}
