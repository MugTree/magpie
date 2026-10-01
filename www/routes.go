package www

import (
	"fmt"
	"net/http"
	"strconv"

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

	// need to list feeds and articles
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("NOT IMPL"))

	})

	// not sure I need this as will all be on hp maybe add later
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

		r.Put("/like/{value}", func(w http.ResponseWriter, r *http.Request) {

			ctx := r.Context()

			articleID, ok := _httpRequireIDParam(w, r, "id")
			if !ok {
				return
			}

			likeValue, err := strconv.Atoi(r.PathValue("value"))
			if err != nil {
				_httpLogAndError(w, r, err.Error())
				return
			}

			if likeValue < 0 && likeValue > 3 {
				_httpLogAndError(w, r, fmt.Sprintf("incorrect like value: %v, needs to be between 0 and 3", likeValue))
				return
			}

			articleLike, err := buildArticleLike(ctx, queries, int64(likeValue), articleID)
			if err != nil {
				_httpLogAndError(w, r, err.Error())
				return
			}

			sse := datastar.NewSSE(w, r)
			sse.PatchElementGostar(articleLike)

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
