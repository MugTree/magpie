package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/mugtree/magpie/db"

	_ "embed"

	_ "github.com/mattn/go-sqlite3"
	"github.com/mugtree/magpie/www"
)

func main() {

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	if err := run(ctx); err != nil {
		www.LogError(err.Error())
	}
}

func run(parent context.Context) error {

	mustEnv := func(key string) string {
		val, ok := os.LookupEnv(key)
		if !ok {
			log.Fatalf("missing .env: %s", key)
		}
		return val
	}

	appPort := mustEnv("APP_PORT")
	appDb := mustEnv("APP_DB")
	appUser := mustEnv("APP_USER")
	appPassword := mustEnv("APP_PASSWORD")

	dbHandle, err := sql.Open("sqlite3", appDb)
	if err != nil {
		return err
	}
	dbHandle.SetMaxOpenConns(1)
	dbHandle.SetMaxIdleConns(1)

	_, _ = dbHandle.Exec(`PRAGMA journal_mode=WAL; PRAGMA busy_timeout=5000;`)

	if err := dbHandle.Ping(); err != nil {
		return err
	}

	queries := db.New(dbHandle)

	webserver := &http.Server{
		Addr:    ":" + appPort,
		Handler: www.SetupHTTPServer(queries, appUser, appPassword),
	}

	application := www.NewApp(dbHandle, queries, webserver)

	// bind OS signal context → app shutdown
	go func() {
		<-parent.Done()
		application.Stop()
	}()

	application.Start()

	if err := application.AppWait(); err != nil {
		return err
	}

	www.LogInfo("server stopped cleanly")
	return nil

}
