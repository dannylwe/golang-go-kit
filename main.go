package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/danny/service/account"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/jinzhu/gorm"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

// const dbSource = APP_DB_URI
func main() {
	fmt.Println("starting application")

	var httpAddr = flag.String("http", ":8080", "http listen address")
	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(os.Stderr)
		logger = log.NewSyncLogger(logger)
		logger = log.With(logger,
			"service", "account",
			"time:", log.DefaultTimestampUTC,
			"caller", log.DefaultCaller,
		)
	}

	level.Info(logger).Log("msg", "service started")
	defer level.Info(logger).Log("msg", "service ended")

	err := godotenv.Load()
	if err != nil {
		level.Error(logger).Log("msg", "Unable to load .env file")
		os.Exit(-1)
	}

	dbSource := os.Getenv("APP_DB_URI")
	fmt.Println(dbSource)

	conn, err := gorm.Open("postgres", dbSource)
	defer conn.Close()

	if err != nil {
		fmt.Print(err)
	}
	conn.Debug().AutoMigrate(&account.Users1{})

	var db *sql.DB
	{
		var err error
		db, err = sql.Open("postgres", dbSource)
		if err != nil {
			level.Error(logger).Log("exit", err)
			os.Exit(-1)
		}
	}

	flag.Parse()
	ctx := context.Background()

	// create new repository
	respository := account.NewRepo(db, logger)
	srv := account.NewService(respository, logger)

	// table := ensureTableExists(db, ctx)

	errs := make(chan error)

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		// clearTable(db, ctx)
		errs <- fmt.Errorf("%s", <-c)
	}()

	endpoints := account.MakeEndpoints(srv)

	go func() {
		fmt.Println("listening on port", *httpAddr)
		handler := account.NewHTTPServer(ctx, endpoints)
		errs <- http.ListenAndServe(*httpAddr, handler)
	}()

	level.Error(logger).Log("exit", <-errs)
}

// database operations
// const tableCreationQuery = `CREATE TABLE IF NOT EXISTS users
// (
// 	id SERIAL,
// 	email TEXT NOT NULL,
// 	password TEXT NOT NULL,
// 	CONSTRAINT Users_pkey PRIMARY KEY (id)
// )`

// func ensureTableExists(db *sql.DB, ctx context.Context) (error) {
// 	if _, err := db.Exec(tableCreationQuery); err != nil {
// 		return err
// 	}
// 	return nil
// }

// func clearTable(db *sql.DB, ctx context.Context) {
// 	db.Exec("DELETE FROM Users")
// 	db.Exec("ALTER SEQUENCE Users_id_seq RESTART WITH 1")
// }
