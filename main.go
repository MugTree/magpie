package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
	"github.com/mugtree/magpie/db"
	"github.com/mugtree/magpie/shared"
)

func main() {

	ctx := context.Background()

	mustEnv := func(key string) string {
		val, ok := os.LookupEnv(key)
		if !ok {
			log.Fatalf("missing .env: %s", key)
		}
		return val
	}

	appDb := mustEnv("APP_DB")

	sqlDB, err := sql.Open("sqlite3", appDb)
	if err != nil {
		shared.LogError(err)
		return
	}

	queries := db.New(sqlDB)

	newArticles, err := shared.AddFeedUpdates(queries, ctx)
	if err != nil {
		shared.LogError(err)
		return
	}

	fmt.Printf("articles created: %v", newArticles)

	_, err = queries.InsertScraperRan(ctx,
		db.InsertScraperRanParams{
			RunType:         "daily",
			ArticlesCreated: newArticles,
		})
	if err != nil {
		shared.LogError(err)
		return
	}

	// it would copy a markdowified version of the article over to another directory
	// when a copy is made this program must not be able to copy again
	// within the copy routine - we need tp have a way to say if a path exists DO NOT OVERWRITE
	// that should be expressed in code and not just be a has been copied flag in the db. Assumption is that everything is always copied

}
