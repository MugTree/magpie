package main

import (
	"bufio"
	"context"
	"database/sql"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/mmcdole/gofeed"

	"github.com/mugtree/magpie/db"
	"github.com/mugtree/magpie/shared"

	_ "github.com/mattn/go-sqlite3"
)

/*
go run ./main.go --urls=./app/db/seed/seed.csv --db=./feeds.db
*/
func main() {

	ctx := context.Background()

	// parse flags
	filePtr := flag.String("urls", "", "the file to get the urls from - needs to be broken over lines")
	dbPtr := flag.String("db", "", "path to the db")

	flag.Parse()
	fmt.Println("urls:", *filePtr)

	// open db
	dbhandle, err := sql.Open("sqlite3", *dbPtr)
	if err != nil {
		shared.LogError(err)
		return
	}
	defer dbhandle.Close()
	queries := db.New(dbhandle)

	// read data from csv to use in program
	f, err := os.Open(*filePtr)
	if err != nil {
		shared.LogError(err)
		return
	}
	defer f.Close()

	seedData := []db.Feed{}

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		parts := strings.Split(scanner.Text(), ",")

		fi := db.Feed{
			Url:                    parts[0],
			CssSelContainer:        parts[1],
			CssSelStart:            parts[2],
			CssSelStop:             parts[3],
			HtmlExtractionStrategy: parts[4],
		}

		seedData = append(seedData, fi)
	}

	if err := scanner.Err(); err != nil {
		shared.LogError(err)
		return
	}

	parser := gofeed.NewParser()

	for _, sd := range seedData {

		goFeed, err := parser.ParseURL(sd.Url)
		if err != nil {
			shared.LogError(err)
			return
		}

		insertedFeed, err := queries.InsertFeed(ctx, db.InsertFeedParams{
			Url:                    goFeed.Link,
			Title:                  goFeed.Title,
			CssSelContainer:        sd.CssSelContainer, //fi.CSSSelectorContainer},
			CssSelStart:            sd.CssSelStart,
			CssSelStop:             sd.CssSelStop,
			HtmlExtractionStrategy: sd.HtmlExtractionStrategy,
		})
		if err != nil {
			shared.LogError(err)
			return
		}

		var createdArticlesCount = 0

		for _, v := range goFeed.Items {

			count, err := shared.AddOrIgnoreArticle(queries, ctx, v, insertedFeed)
			if err != nil {
				shared.LogError(err)
			}

			createdArticlesCount = createdArticlesCount + count
		}

		_, err = queries.InsertScraperRan(ctx, db.InsertScraperRanParams{RunType: "seed", ArticlesCreated: int64(createdArticlesCount)})

	}
}
