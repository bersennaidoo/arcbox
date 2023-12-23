package dbc

import (
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql"
	"github.com/spf13/viper"
)

func New(config *viper.Viper) *sql.DB {
	connectionString := config.GetString("database.connection_string")
	maxIdleConnections := config.GetInt("database.max_idle_connections")
	maxOpenConnections := config.GetInt("database.max_open_connections")
	connectionMaxLifetime := config.GetDuration("database.connection_max_lifetime")
	driverName := config.GetString("database.driver_name")

	if connectionString == "" {
		log.Fatalf("Database connectin string is missing")
	}

	dbClient, err := sql.Open(driverName, connectionString)
	if err != nil {
		log.Fatalf("Error while initializing database: %v", err)
	}

	dbClient.SetMaxIdleConns(maxIdleConnections)
	dbClient.SetMaxOpenConns(maxOpenConnections)
	dbClient.SetConnMaxLifetime(connectionMaxLifetime)

	err = dbClient.Ping()
	if err != nil {
		dbClient.Close()
		log.Fatalf("Error while validating database: %v", err)
	}

	log.Println("Mysql Connected")

	return dbClient
}
