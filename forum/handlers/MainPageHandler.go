package handlers

import (
	"errors"
	"fmt"
	"net/http"

	"forum/forum/domain"
	"forum/forum/internal"
)

func (hh *HttpHandler) MainHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {

		posts, err := hh.business.GetAllPosts()
		if err != nil {
			fmt.Println("Cant get Posts")
			return
		}
		username, err := hh.GetUsername(w, r)
		if err != nil {
			if errors.Is(err, domain.ErrSessionNotFound) {
				internal.RenderMainPage(w, r, username, posts)
				return
			}
			internal.RenderMainPage(w, r, username, posts)
			return
		}
		if username.Username == "" {
			internal.RenderMainPage(w, r, username, posts)
			return
		}

		internal.RenderMainPage(w, r, username, posts)
	}
}
