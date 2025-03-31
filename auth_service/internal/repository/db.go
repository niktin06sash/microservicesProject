package repository

import (
	"auth_service/configs"
	"context"
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"

	"github.com/redis/go-redis/v9"
)

type DBInterface interface {
	Open(driverName, connectionString string) (*sql.DB, error)
	Ping(db *sql.DB) error
	Close(db *sql.DB)
	SetConfig(cfg DBConfig)
}

type DBConfig struct {
	Driver   string
	Host     string
	Port     int
	User     string
	Password string
	Name     string
	SSLMode  string
}

type DBObject struct {
	dbConfig DBConfig
}

func (d *DBObject) SetConfig(cfg DBConfig) {
	d.dbConfig = cfg
}

func (d *DBObject) Open(driverName string, connectionString string) (*sql.DB, error) {
	db, err := sql.Open(driverName, connectionString)
	if err != nil {
		log.Printf("Sql-Open error %v", err)
		return nil, err
	}
	return db, nil
}

func (d *DBObject) Ping(db *sql.DB) error {
	err := db.Ping()
	if err != nil {
		log.Printf("Sql-Ping error %v", err)
		return err
	}
	return nil
}

func (d *DBObject) Close(db *sql.DB) {
	err := db.Close()
	if err != nil {
		log.Printf("Sql-Close error %v", err)
	}
}

func ConnectToDb(cfg configs.Config) (*sql.DB, DBInterface, error) {
	dbInterface := &DBObject{}

	dbConfig := DBConfig{
		Driver:   cfg.Database.Driver,
		Host:     cfg.Database.Host,
		Port:     cfg.Database.Port,
		User:     cfg.Database.User,
		Password: cfg.Database.Password,
		Name:     cfg.Database.Name,
		SSLMode:  cfg.Database.SSLMode,
	}
	dbInterface.SetConfig(dbConfig)

	connectionString := fmt.Sprintf("postgresql://%s:%s@%s:%d/%s?sslmode=%s",
		dbConfig.User, dbConfig.Password, dbConfig.Host, dbConfig.Port, dbConfig.Name, dbConfig.SSLMode)

	db, err := dbInterface.Open(dbConfig.Driver, connectionString)
	if err != nil {
		log.Printf("Sql-Open error %v", err)
		return nil, nil, err
	}
	err = dbInterface.Ping(db)
	if err != nil {
		dbInterface.Close(db)
		log.Printf("Sql-Ping error %v", err)
		return nil, nil, err
	}
	return db, dbInterface, nil
}

type RedisInterface interface {
	Open(host string, port int, password string, db int) (*redis.Client, error)
	Ping(client *redis.Client) error
	Close(client *redis.Client) error
}

type RedisObject struct{}

func (r *RedisObject) Open(host string, port int, password string, db int) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", host, port),
		Password: password,
		DB:       db,
	})
	return client, nil
}

func (r *RedisObject) Ping(client *redis.Client) error {
	_, err := client.Ping(context.Background()).Result()
	return err
}

func (r *RedisObject) Close(client *redis.Client) error {
	return client.Close()
}

func ConnectToRedis(cfg configs.Config) (*redis.Client, RedisInterface, error) {
	redisInterface := &RedisObject{}

	client, err := redisInterface.Open(cfg.Redis.Host, cfg.Redis.Port, cfg.Redis.Password, cfg.Redis.DB)
	if err != nil {
		log.Printf("Sql-Open error %v", err)
		return nil, nil, err
	}
	err = redisInterface.Ping(client)
	if err != nil {
		redisInterface.Close(client)
		log.Printf("Sql-Ping error %v", err)
		return nil, nil, err
	}
	return client, redisInterface, nil
}
