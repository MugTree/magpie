package www

import (
	"net/http"

	"github.com/mugtree/magpie/db"

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

func articleHandler(_ *db.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}

func articleUpdateHandler(_ *db.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}
