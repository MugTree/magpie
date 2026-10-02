package www

import (
	"fmt"
	"strconv"

	"github.com/mugtree/magpie/db"
	. "maragu.dev/gomponents"
	ds "maragu.dev/gomponents-datastar"
	. "maragu.dev/gomponents/components"
	. "maragu.dev/gomponents/html"
)

type pageProps struct {
	Title       string
	Description string
}

func layout(props pageProps, children ...Node) Node {
	// Hash the paths for easy cache busting on changes

	return HTML5(HTML5Props{
		Title:       props.Title,
		Description: props.Description,
		Language:    "en",
		Head: []Node{
			Script(Src("/public/js/datastar.js"), Type("module")),
			// Link(Rel("stylesheet"), Href("/public/css/app.css")),
			Link(Rel("stylesheet"), Href("/public/css/main.css")),
		},
		Body: Group(children),
	})

}

func articlePage(article db.Article, html string, sigs map[string]any) Node {

	return Div(
		H1(Text(article.Title)),
		P(Text(fmt.Sprint(article.Link, " ", article.Published))),
		articleLike(article.ID, article.Starred),
		Hr(),
		Main(
			ID("article"),
			Class(`article`),
			ds.Signals(sigs),

			// stars

			// server rendered HTML
			authorHTML(html),

			// markdown
			Div(ID("markdown"),
				Textarea(ID("editor"),
					ds.Bind("edit"),
					ds.On("input", fmt.Sprintf("@get('/article/%v/write')", article.ID)),
					Text(article.Markdown),
				),
				Button(
					ds.On("click", fmt.Sprintf("@patch('/article/%v/save')", article.ID)),
					Text("submit"),
				),
			),
		))
}

func authorHTML(html string) Node {
	return Div(ID("html"),
		Raw(html),
	)
}

func articleLike(articleID int64, starsValue int64) Node {
	return Div(ID("star-value-bottom"), ds.On("click", fmt.Sprintf("@put('/article/%v/like/%v')", articleID, starsValue)),

		Text("Starred:"),

		Img(
			Width("60px"),
			Src(fmt.Sprintf("/public/img/%v-star.png", starsValue)),
		),
	)
}

func homePage(summaries []feedSummary) Node {

	feedPagination := func(fsm feedSummary) Node {
		links := []Node{}
		for i := range fsm.LinksRequired {
			pageNumber := i + 1
			links = append(links,
				A(
					ds.On("click", fmt.Sprintf("@get('/feed/%v/page/%v')", fsm.FeedID, pageNumber)),
					Text(strconv.FormatInt(pageNumber, 10)),
					Classes{"link": true, "underline": fsm.PageID == pageNumber},
				),
			)
		}
		return Group(links)
	}

	return Div(ID("homepage"),
		Div(ID("feeds"),
			Map(summaries, func(fs feedSummary) Node {
				return Div(ID(fmt.Sprintf("feed-box-%v", fs.FeedID)),
					Class("feedbox"),
					H2(Text(fs.Name), ds.On("click", fmt.Sprintf("@get('/feed/%v/page/1')", fs.FeedID))),

					If(len(fs.Articles) > 0,
						Div(
							Map(fs.Articles,
								func(a db.SelectArticlesByFeedIDWithLimitRow) Node {

									return H3(A(Text(a.ArticleTitle), Href(fmt.Sprintf("/article/%v/read", a.ArticleID))))
								},
							),
							Div(feedPagination(fs)),
						),
					),
				)

			}),
		),
	)

}
