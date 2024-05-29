// package main

// import (
// 	"database/sql"
// 	"dreampicai/db"
// 	"log"
// 	"os"

// 	"github.com/golang-migrate/migrate/v4"
// 	"github.com/golang-migrate/migrate/v4/database/postgres"
// 	_ "github.com/golang-migrate/migrate/v4/source/file"
// 	"github.com/joho/godotenv"
// )

// func createDB() (*sql.DB, error) {
// 	if err := godotenv.Load(); err != nil {
// 		return nil, err
// 	}

// 	var (
// 		host   = os.Getenv("DB_HOST")
// 		user   = os.Getenv("DB_USER")
// 		pass   = os.Getenv("DB_PASSWORD")
// 		dbName = os.Getenv("DB_NAME")
// 	)
// 	return db.CreateDatabase(dbName, user, pass, host)
// }

// func main() {
// 	db, err := createDB()
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	// create migration instance
// 	driver, err := postgres.WithInstance(db, &postgres.Config{})
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	// post to your migration files. Here we're using local files, but it could be others source
// 	m, err := migrate.NewWithDatabaseInstance(
// 		"file://cmd/migrate/migrations",
// 		"postgres",
// 		driver,
// 	)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	cmd := os.Args[len(os.Args)-1]
// 	if cmd == "up" {
// 		if err := m.Up(); err != nil && err != migrate.ErrNoChange {
// 			log.Fatal(err)
// 		}
// 	}
// 	if cmd == "down" {
// 		if err := m.Down(); err != nil && err != migrate.ErrNoChange {
// 			log.Fatal(err)
// 		}
// 	}

// 	defer m.Close()

// }
