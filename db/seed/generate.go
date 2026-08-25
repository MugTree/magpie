package main

import (
	"bufio"
	"context"
	"database/sql"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
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

	mustEnv := func(key string) string {
		val, ok := os.LookupEnv(key)
		if !ok {
			log.Fatalf("missing .env: %s", key)
		}
		return val
	}

	markdownDir := mustEnv("MARKDOWN_DIR")

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
			FolderName:             parts[1],
			CssSelContainer:        parts[2],
			CssSelStart:            parts[3],
			CssSelStop:             parts[4],
			HtmlExtractionStrategy: parts[5],
		}

		seedData = append(seedData, fi)
	}

	if err := scanner.Err(); err != nil {
		shared.LogError(err)
		return
	}

	parser := gofeed.NewParser()

	for _, d := range seedData {

		goFeed, err := parser.ParseURL(d.Url)
		if err != nil {
			shared.LogError(err)
			return
		}

		feed, err := queries.InsertFeed(ctx,
			db.InsertFeedParams{
				Url:                    goFeed.Link,
				Title:                  goFeed.Title,
				CssSelContainer:        d.CssSelContainer, //fi.CSSSelectorContainer},
				CssSelStart:            d.CssSelStart,
				CssSelStop:             d.CssSelStop,
				FolderName:             d.FolderName,
				HtmlExtractionStrategy: d.HtmlExtractionStrategy,
			})
		if err != nil {
			shared.LogError(err)
			return
		}

		markdownPath := filepath.Join(markdownDir, d.FolderName)
		err = os.MkdirAll(markdownPath, 0755)
		if err != nil {
			shared.LogError(err)
			return
		}

		var createdArticlesTally = 0

		for _, v := range goFeed.Items {

			article, err := shared.InsertOrIgnoreArticle(queries, ctx, v, feed)
			if err != nil {
				shared.LogError(err)
			}

			if article.ID != 0 {

				err := shared.CreateMarkdown(article, markdownPath)
				if err != nil {
					shared.LogError(err)
					return
				}

				createdArticlesTally = createdArticlesTally + 1
			}

		}

		_, err = queries.InsertScraperRan(
			ctx,
			db.InsertScraperRanParams{
				RunType:         "seed",
				ArticlesCreated: int64(createdArticlesTally)},
		)

	}
}
