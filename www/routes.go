package www

import (
	"net/http"

	"github.com/goforj/godump"
	"github.com/mugtree/magpie/db"
	"github.com/starfederation/datastar-go/datastar"

	"github.com/go-chi/chi/v5"
	// "github.com/goforj/godump"
)

func getRouter(queries *db.Queries) chi.Router {
	r := chi.NewRouter()
	// r.Use(httpDebugRequest)
	r.Handle("/public/*", _httpNeuterDirectory(http.FileServer(http.FS(staticFS))))

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("NOT IMPL"))

	})
	r.Get("/feed/{id}", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("NOT IMPL"))

	})
	r.Route("/article/{id}", func(r chi.Router) {

		r.Get("/read", func(w http.ResponseWriter, r *http.Request) {

			id, ok := _httpRequireIDParam(w, r, "id")
			if !ok {
				return
			}

			articlePage, article, err := buildArticlePage(r.Context(), queries, id)
			if err != nil {
				_httpLogAndError(w, r, err.Error())
				return
			}

			layout(
				pageProps{Title: article.Title},
				articlePage,
			).Render(w)

		})

		r.Get("/write", func(w http.ResponseWriter, r *http.Request) {

			as, ok := _signalsOrError(w, r)
			if !ok {
				return
			}

			notatedArticle := notateArticlePage(as)

			sse := datastar.NewSSE(w, r)
			sse.PatchElementGostar(notatedArticle)

		})

		r.Patch("/save", func(w http.ResponseWriter, r *http.Request) {

			as, ok := _signalsOrError(w, r)
			if !ok {
				return
			}

			updatedPage, updatedSignals, err := updateArticlePage(r.Context(), queries, as)
			if err != nil {
				_httpLogAndError(w, r, err.Error())
				return
			}

			sse := datastar.NewSSE(w, r)
			sse.PatchElementGostar(updatedPage)
			sse.MarshalAndPatchSignals(updatedSignals)
		})

	})
	return r
}

func _signalsOrError(w http.ResponseWriter, r *http.Request) (articleSignals, bool) {
	as := articleSignals{}
	if err := datastar.ReadSignals(r, &as); err != nil {
		_httpLogAndError(w, r, err.Error())
		return as, false
	}

	return as, true
}

func DummyCallToPreserveGoDump() {
	godump.Dump("")
}
