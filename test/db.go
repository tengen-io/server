package test

import (
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	"github.com/lib/pq"
	"github.com/tengen-io/server/db"
	"golang.org/x/crypto/bcrypt"
	"log"
	"os"
	"strconv"
	"testing"
	"time"
)

const (
	Password = "hunter2"
)

var dbHandle *sqlx.DB

func DB() *sqlx.DB {
	return dbHandle
}

func TestMain(m *testing.M, suite string) {
	godotenv.Load("../.env.test")

	db := initializeTestDB(suite)

	setupFixtures(db)
//	defer teardownFixtures(db)
	rv := m.Run()

	os.Exit(rv)

}

func initializeTestDB(name string) *sqlx.DB {
	port, err := strconv.Atoi(os.Getenv("TENGEN_DB_PORT"))
	if err != nil {
		log.Fatal("Could not parse TENGEN_DB_PORT: ", err)
	}
	database := os.Getenv("TENGEN_DB_DATABASE") + "_" + name
	rootConfig := db.PostgresDBConfig{
		Host:     os.Getenv("TENGEN_DB_HOST"),
		Port:     port,
		User:     os.Getenv("TENGEN_DB_USER"),
		Database: "postgres",
		Password: os.Getenv("TENGEN_DB_PASSWORD"),
	}

	rootDb, err := db.NewPostgresDb(rootConfig)
	if err != nil {
		log.Fatal("unable to connect to DB.", err)
	}

	_, err = rootDb.Exec("DROP DATABASE IF EXISTS " + database)
	if err != nil {
		log.Fatalf("unable to drop existing test db: %s", err)
	}

	_, err = rootDb.Exec("CREATE DATABASE " + database)
	if err != nil {
		log.Fatalf("unable to create test db: %s", err)
	}

	var config = rootConfig
	config.Database = database
	rv, err := db.NewPostgresDb(config)
	if err != nil {
		log.Fatalf("unable to create connection to test db: %s", err)
	}

	dbHandle = rv

	return rv
}

func setupFixtures(db *sqlx.DB) {
	driver, err := postgres.WithInstance(db.DB, &postgres.Config{})
	if err != nil {
		panic(err)
	}
	fsrc, err := (&file.File{}).Open("file://../db/migrations")
	if err != nil {
		panic(err)
	}
	m, err := migrate.NewWithInstance("file", fsrc, "postgres", driver)
	if err != nil {
		panic(err)
	}
	m.Down()
	err = m.Up()
	if err != nil {
		panic(err)
	}

	writeFixtures(db)
}

func writeFixtures(db *sqlx.DB) {
	hash, _ := bcrypt.GenerateFromPassword([]byte(Password), 4)
	now := pq.FormatTimestamp(time.Now().UTC())
	res := db.QueryRow("INSERT INTO identities (email, password_hash, created_at, updated_at) VALUES ($1, $2, $3, $4) RETURNING id", "test1@tengen.io", hash, now, now)
	var id int64
	res.Scan(&id)
	db.MustExec("INSERT INTO users (identity_id, name, created_at, updated_at) VALUES ($1, $2, $3, $4)", id, "Test User 1", now, now)

	db.MustExec("INSERT INTO games (type, state, board_size, created_at, updated_at) VALUES ($1, $2, $3, $4, $5)", "STANDARD", "NEGOTIATION", 19, now, now)
	db.MustExec("INSERT INTO games (type, state, board_size, created_at, updated_at) VALUES ($1, $2, $3, $4, $5)", "STANDARD", "IN_PROGRESS", 19, now, now)
}

func teardownFixtures(db *sqlx.DB) {
	driver, err := postgres.WithInstance(db.DB, &postgres.Config{})
	if err != nil {
		panic(err)
	}
	m, err := migrate.NewWithDatabaseInstance("file://db/migrations", "postgres", driver)
	if err != nil {
		panic(err)
	}
	m.Drop()
}

