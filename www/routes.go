package www

import (
	"net/http"

	"github.com/mugtree/magpie/db"

	"github.com/go-chi/chi/v5"
	// "github.com/goforj/godump"
)

func getRouter(_ *db.Queries) chi.Router {
	r := chi.NewRouter()
	//r.Use(httpDebugRequest)
	r.Handle("/public/*", httpNeuterDirectory(http.FileServer(http.FS(staticFS))))
	return r
}
