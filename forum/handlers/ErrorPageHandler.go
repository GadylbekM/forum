package handlers

import (
	"net/http"

	"forum/forum/internal"
)

func (hh *HttpHandler) Handle404(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(404)
	internal.RenderErrorPage(w, r)
}

func (hh *HttpHandler) Handle403(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(403)
	internal.RenderError403Page(w, r)
}
