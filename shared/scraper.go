package shared

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/mugtree/magpie/db"

	"github.com/gocolly/colly/v2"
	"github.com/mmcdole/gofeed"

	htmltomarkdown "github.com/JohannesKaufmann/html-to-markdown/v2"
)

func InsertFeedUpdates(queries *db.Queries, ctx context.Context, markdownPath string) (int64, error) {

	feeds, err := queries.SelectAllFeeds(ctx)
	if err != nil {
		return 0, fmt.Errorf("get feeds: %w", err)
	}

	parser := gofeed.NewParser()
	parser.Client = &http.Client{
		Timeout: 10 * time.Second,
	}

	var createdArticlesTally = 0

	for _, feed := range feeds {

		goFeed, err := parser.ParseURL(fmt.Sprintf("%s/feed/", feed.Url))
		if err != nil {
			return 0, fmt.Errorf("parse feed %s: %w", feed.Url, err)
		}

		// feed may be missing or down
		if goFeed == nil {
			continue
		}

		markdownPath := filepath.Join(markdownPath, feed.FolderName)

		for _, v := range goFeed.Items {

			article, err := InsertOrIgnoreArticle(queries, ctx, v, feed)
			if err != nil {
				LogError(err)
			}

			if article.ID != 0 {

				err := CreateMarkdown(article, feed, markdownPath)
				if err != nil {
					return 0, fmt.Errorf("parse feed %s: %w", feed.Url, err)
				}

				createdArticlesTally = createdArticlesTally + 1
			}

		}
	}

	_, err = queries.InsertScraperRan(ctx, db.InsertScraperRanParams{
		RunType:         "daily",
		ArticlesCreated: int64(createdArticlesTally),
	})

	return int64(createdArticlesTally), nil
}

func InsertOrIgnoreArticle(queries *db.Queries, ctx context.Context, item *gofeed.Item, feed db.Feed) (db.Article, error) {

	var article db.Article

	publishedDate := feedItemDate(item)
	dateFound := time.Now()

	scrapeHTML := func(ep PageScrapeParams) (string, error) {

		pageHtmlContent := ""

		c := colly.NewCollector()

		c.OnHTML(ep.Container, func(h *colly.HTMLElement) {
			pageHtmlContent = ExtractHTMLRangeFlat(h.DOM, ep.ClipStartPoint, ep.ClipEndPoint)
		})

		if err := c.Visit(ep.Link); err != nil {
			return "", fmt.Errorf("error using colly to visit page: %v - %v", ep.Link, err)
		}

		return pageHtmlContent, nil

	}

	html, err := scrapeHTML(PageScrapeParams{
		Link:           item.Link,
		Container:      feed.CssSelContainer,
		ClipStartPoint: feed.CssSelStart,
		ClipEndPoint:   feed.CssSelStop,
	})
	if err != nil {
		return article, fmt.Errorf("error getting site html: %v", err)
	}

	processed, err := ProcessScrapedHTML(html)
	if err != nil {
		return article, fmt.Errorf("error running ProcessScrapedHTML: %v", err)
	}

	// create the markdown here ...

	err = queries.InsertOrIgnoreArticle(ctx, db.InsertOrIgnoreArticleParams{
		FeedID:        feed.ID,
		Title:         item.Title,
		Link:          item.Link,
		Published:     publishedDate,
		DateFound:     &dateFound,
		Summary:       item.Description,
		ScrapedHtml:   html,
		FormattedHtml: processed,
	})
	if err != nil {
		return article, fmt.Errorf("error inserting article: %v", err)
	}

	newArticle, err := queries.SelectArticleByFeedIDAndLink(ctx,
		db.SelectArticleByFeedIDAndLinkParams{
			FeedID: feed.ID,
			Link:   item.Link,
		})

	if err != nil {
		if err == sql.ErrNoRows {
			return article, nil
		}
		return article, fmt.Errorf("error selecting article: %v", err)
	}

	return newArticle, nil

}

func CreateMarkdown(article db.Article, feed db.Feed, markdownPath string) error {

	slug := func(s string) string {
		s = strings.ToLower(s)
		re := regexp.MustCompile(`[^a-z0-9]+`)
		s = re.ReplaceAllString(s, "-")
		return strings.Trim(s, "-")
	}

	date := article.DateFound.Format("2006-01-02")

	fileName := date + "-" + slug(article.Title) + ".md"

	f, err := os.Create(filepath.Join(markdownPath, fileName))
	if err != nil {
		return err
	}
	defer f.Close()

	content, err := htmltomarkdown.ConvertString(article.FormattedHtml)
	if err != nil {
		return err
	}

	tags, err := htmltomarkdown.ConvertString(fmt.Sprintf(`<p>#annotations #%v</p>`, strings.ReplaceAll(feed.FolderName, "_", "")))
	if err != nil {
		return err
	}

	link, err := htmltomarkdown.ConvertString(fmt.Sprintf(`<p><a href="%v">%v</a> | %v</p>`, feed.Url, feed.Title, date))
	if err != nil {
		return err
	}

	content = fmt.Sprintf("\n\n%v\n\n## %v\n\n%v\n\n%v", tags, article.Title, link, content)

	_, err = f.Write([]byte(content))

	if err != nil {
		LogError(err)
		return err
	}

	return nil

}

func feedItemDate(item *gofeed.Item) *time.Time {
	if item.PublishedParsed != nil {
		return item.PublishedParsed
	}

	if item.UpdatedParsed != nil {
		return item.UpdatedParsed
	}

	return nil
}
