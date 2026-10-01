package www

import (
	"net/http"

	"github.com/mugtree/magpie/db"
	"github.com/russross/blackfriday/v2"

	"github.com/go-chi/chi/v5"
	// "github.com/goforj/godump"

	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

func getRouter(queries *db.Queries) chi.Router {
	r := chi.NewRouter()
	// r.Use(httpDebugRequest)
	r.Handle("/public/*", httpNeuterDirectory(http.FileServer(http.FS(staticFS))))

	r.Get("/", homeHandler(queries))
	r.Get("/feed/{id}", feedHandler(queries))
	r.Get("/article/{id}", articleHandler(queries))
	r.Post("/article", articleUpdateHandler(queries))
	return r
}

func homeHandler(_ *db.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		layout(pageProps{Title: "Home"}, Div(Text("Home!"))).Render(w)
	}
}

func feedHandler(_ *db.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}

func articleHandler(queries *db.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		id, ok := httpRequireIDParam(w, r, "id")
		if !ok {
			return
		}

		ctx := r.Context()
		article, err := queries.SelectArticleByID(ctx, id)
		if err != nil {
			httpLogAndError(w, r, err.Error())
			return
		}

		md := blackfriday.Run([]byte(article.Markdown))

		page := Main(
			Div(
				Textarea(ID("editor"),
					Text(article.Markdown)),
			),
			Div(
				Raw(string(md)),
			),
		)

		layout(pageProps{Title: article.Title}, page).Render(w)

	}
}

func articleUpdateHandler(_ *db.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}
