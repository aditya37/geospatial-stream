package service

import (
	"net/http"

	"github.com/gorilla/mux"
)

type httpHandler struct {
	muxRouter *mux.Router
}

func NewHttpHandler() *httpHandler {
	router := mux.NewRouter()
	return &httpHandler{
		muxRouter: router,
	}
}

func (h *httpHandler) Handler() http.Handler {
	// health check
	h.muxRouter.Methods(http.MethodGet).Path("/").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("{'status':'up'}"))
	})
	return h.muxRouter
}
