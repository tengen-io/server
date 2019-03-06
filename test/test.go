package test

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	"github.com/lib/pq"
	_ "github.com/lib/pq"
	"github.com/tengen-io/server/db"
	"golang.org/x/crypto/bcrypt"
	"log"
	"os"
	"strconv"
	"testing"
	"time"
)

var testDb string

const (
	Password = "hunter2"
)

func MakeDb() *sqlx.DB {
	godotenv.Load("../.env.test")
	port, err := strconv.Atoi(os.Getenv("TENGEN_DB_PORT"))
	if err != nil {
		log.Fatal("Could not parse TENGEN_DB_PORT: ", err)
	}

	dbName := os.Getenv("TENGEN_DB_DATABASE")
	if testDb != "" {
		dbName += "_" + testDb
	}

	config := &db.PostgresDBConfig{
		Host:     os.Getenv("TENGEN_DB_HOST"),
		Port:     port,
		User:     os.Getenv("TENGEN_DB_USER"),
		Database: dbName,
		Password: os.Getenv("TENGEN_DB_PASSWORD"),
	}

	db, err := db.NewPostgresDb(config)
	if err != nil {
		log.Fatal("Unable to connect to DB.", err)
	}

	return db
}

func Main(m *testing.M) {
	fmt.Println("setting up")
	SetupFixtures()
	defer TeardownFixtures()
	rv := m.Run()
	fmt.Println("tearing down")
	os.Exit(rv)
}

// TODO(eac) this whole thing is horrible
func SetupFixtures() {
	var bytes [6]byte
	rand.Read(bytes[:])
	suffix := base64.StdEncoding.EncodeToString(bytes[:])
	db := MakeDb()
	dbName := os.Getenv("TENGEN_DB_DATABASE") + "_" + suffix
	db.MustExec("CREATE DATABASE "+dbName)
	testDb = dbName
	db = MakeDb()
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
	err = m.Up()
	if err != nil {
		panic(err)
	}
	fixtures()
}

func TeardownFixtures() {
	db := MakeDb()
	driver, err := postgres.WithInstance(db.DB, &postgres.Config{})
	if err != nil {
		panic(err)
	}
	m, err := migrate.NewWithDatabaseInstance("file://../db/migrations", "postgres", driver)
	if err != nil {
		panic(err)
	}
	m.Down()
}

func fixtures() {
	db := MakeDb()
	hash, _ := bcrypt.GenerateFromPassword([]byte(Password), 4)
	now := pq.FormatTimestamp(time.Now().UTC())
	res := db.QueryRow("INSERT INTO identities (email, password_hash, created_at, updated_at) VALUES ($1, $2, $3, $4) RETURNING id", "test1@tengen.io", hash, now, now)
	var id int64
	res.Scan(&id)
	db.MustExec("INSERT INTO users (identity_id, name, created_at, updated_at) VALUES ($1, $2, $3, $4)", id, "Test User 1", now, now)

	db.MustExec("INSERT INTO games (type, state, board_size, created_at, updated_at) VALUES ($1, $2, $3, $4, $5)", "STANDARD", "INVITATION", 19, now, now)
	db.MustExec("INSERT INTO games (type, state, board_size, created_at, updated_at) VALUES ($1, $2, $3, $4, $5)", "STANDARD", "IN_PROGRESS", 19, now, now)
}
