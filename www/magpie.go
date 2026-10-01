package www

import (
	"context"

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

	html := blackfriday.Run([]byte(article.Markdown))

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

	html := blackfriday.Run([]byte(article.Markdown))

	sigs, err := lib.StructToMap(as)
	if err != nil {
		return nil, as, err
	}

	return articlePage(article, string(html), sigs), as, nil

}

func notateArticlePage(as articleSignals) Node {
	note := string(blackfriday.Run([]byte(as.Edit)))
	return htmlLayout(note)
}
