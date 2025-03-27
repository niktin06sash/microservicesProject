package repository

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
	"github.com/spf13/viper"
)

type DBInterface interface {
	Open(driverName, dataSourceName string) (*sql.DB, error)
	Ping(db *sql.DB) error
	Close(db *sql.DB) error
}

type DBObject struct{}

func (d DBObject) Open(driverName, dataSourceName string) (*sql.DB, error) {
	return sql.Open(driverName, dataSourceName)
}

func (d DBObject) Ping(db *sql.DB) error {
	return db.Ping()
}

func (d DBObject) Close(db *sql.DB) error {
	return db.Close()
}
func LoadConfig(path string) {
	viper.SetConfigFile(path)
	viper.SetConfigType("yml")

	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("Error reading config file, %s", err)
		os.Exit(1)
	}

	fmt.Println("Successfully loaded configuration!")
}
func ConnectToDb(dbInterface DBInterface) (*sql.DB, error) {
	dbDriver := viper.GetString("database.driver")
	dbHost := viper.GetString("database.host")
	dbPort := viper.GetInt("database.port")
	dbUser := viper.GetString("database.user")
	dbPass := viper.GetString("database.password")
	dbName := viper.GetString("database.name")
	dbSSLMode := viper.GetString("database.sslmode")
	connectionString := fmt.Sprintf("postgresql://%s:%s@%s:%d/%s?sslmode=%s", dbUser, dbPass, dbHost, dbPort, dbName, dbSSLMode)

	db, err := dbInterface.Open(dbDriver, connectionString)
	if err != nil {
		return nil, err
	}
	err = dbInterface.Ping(db)
	if err != nil {
		dbInterface.Close(db)
		return nil, fmt.Errorf("error pinging database: %w", err)
	}
	return db, nil
}
