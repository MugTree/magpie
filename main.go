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

	/* the idea here is that you simply run the app and id downloads the latest feed entries and adds the feeds to the database */

	sqlDB, err := sql.Open("sqlite3", appDb)
	if err != nil {
		fmt.Printf("error opening the db: %v", err)
		return
	}

	queries := db.New(sqlDB)

	// go to the web and get the feed contents (articles / urls)
	// add the new ones to the database (new is based upon a unique key: feedID and link)
	// typically this would be run once a day

	articlesCreated, err := shared.GetFeedUpdates(queries, ctx)
	if err != nil {
		fmt.Printf("error running GetFeedUpdates: %v", err)
		return
	}

	fmt.Printf("articles created: %v", articlesCreated)

	_, err = queries.InsertScraperRan(ctx, db.InsertScraperRanParams{RunType: "daily", ArticlesCreated: articlesCreated})
	if err != nil {
		fmt.Println(err)
		return
	}

	// it would copy a markdowified version of the article over to another directory

	// when a copy is made this program must not be able to copy again

	// within the copy routine - we need tp have a way to say if a path exists DO NOT OVERWRITE

	// that should be expressed in code and not just be a has been copied flag in the db. Assumption is that everything is always copied

}
