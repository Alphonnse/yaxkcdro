package http

import (
	"net/http"

	"github.com/Alphonnse/yaxkcdro/internal/UseCase/api"
)

type Router struct {
	server *http.Server
	mux    *http.ServeMux
	api    api.APIUseCaseService
}

func NewRouter(server *http.Server, mux *http.ServeMux, api api.APIUseCaseService) {
	r := &Router{
		server: server,
		mux:    mux,
		api:    api,
	}

	r.mux.HandleFunc("/healthz", r.healthz)
	r.mux.HandleFunc("/update", r.update)
}

func (*Router) healthz(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write([]byte("OK"))
}

func (router *Router) update(w http.ResponseWriter, r *http.Request) {
	err := router.api.UpdateComicsCount()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(err.Error()))
		return
	}
	w.WriteHeader(http.StatusOK)
}
