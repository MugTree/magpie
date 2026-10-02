package www

import (
	"context"

	"github.com/microcosm-cc/bluemonday"
	"github.com/mugtree/magpie/db"
	"github.com/mugtree/magpie/lib"
	"github.com/russross/blackfriday/v2"

	. "maragu.dev/gomponents"
)

type articleSignals struct {
	Edit      string `json:"edit"`
	ArticleID int64  `json:"article_id"`
	Saved     bool   `json:"saved"`
}

func buildArticlePage(ctx context.Context, queries *db.Queries, articleID int64) (Node, db.Article, error) {

	article, err := queries.SelectArticleByID(ctx, articleID)
	if err != nil {
		return nil, db.Article{}, err
	}

	as := articleSignals{Edit: article.Markdown, ArticleID: articleID}

	sigs, err := lib.StructToMap(as)
	if err != nil {
		return nil, db.Article{}, err
	}

	html := formatMarkown(article.Markdown)

	return articlePage(article, string(html), sigs), article, nil

}

func updateArticlePage(ctx context.Context, queries *db.Queries, as articleSignals) (Node, articleSignals, error) {

	article, err := queries.UpdateArticleByID(
		ctx,
		db.UpdateArticleByIDParams{ID: as.ArticleID, Markdown: as.Edit},
	)
	if err != nil {
		return nil, as, err
	}
	as.Saved = true

	html := formatMarkown(article.Markdown) //blackfriday.Run([]byte(article.Markdown))

	sigs, err := lib.StructToMap(as)
	if err != nil {
		return nil, as, err
	}

	return articlePage(article, string(html), sigs), as, nil

}

func notateArticlePage(as articleSignals) Node {
	note := formatMarkown(as.Edit)
	return authorHTML(note)
}

func updateArticleLike(ctx context.Context, queries *db.Queries, starredValue int64, articleID int64) (Node, error) {

	updatedValue := func(currentValue int64) int64 {
		if currentValue == 3 {
			return 0
		}
		return currentValue + 1
	}(starredValue)

	article, err := queries.UpdateArticleSetStarredValue(ctx,
		db.UpdateArticleSetStarredValueParams{
			Starred: int64(updatedValue),
			ID:      articleID},
	)
	if err != nil {
		return nil, err
	}

	return articleLike(articleID, article.Starred), nil
}

func buildHomePage(ctx context.Context, queries *db.Queries) ([]Node, error) {

	articles, err := queries.SelectArticlesWithFeedName(ctx)
	if err != nil {
		return nil, err
	}

	articlesMap := map[string][]db.Article{}

	for i, v := range articles {
		key := articles[i].FeedTitle
		articlesMap[key] = append(articlesMap[key], db.Article{ID: v.ID, Title: v.Title})
	}

	boxes := []Node{}
	for k, v := range articlesMap {
		boxes = append(boxes, homePageBox(k, v))
	}

	return boxes, nil

}

// helpers
func formatMarkown(input string) string {
	unsafe := blackfriday.Run([]byte(input))
	return string(bluemonday.UGCPolicy().SanitizeBytes(unsafe))
}
