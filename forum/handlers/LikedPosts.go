package handlers

import (
	"fmt"
	"net/http"

	"forum/forum/internal"
)

func (hh *HttpHandler) HandleLikedPosts(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		sessionCookie, err := r.Cookie("session_id")
		if err != nil {
			fmt.Println(err)
			http.Redirect(w, r, "/access_denied", http.StatusSeeOther)
			return
		}
		sessionID := sessionCookie.Value

		session, err := hh.business.Session(sessionID)
		if err != nil {
			fmt.Println(err)
			http.Redirect(w, r, "/access_denied", http.StatusSeeOther)
			return
		}
		likedPosts, err := hh.business.GetLikedPosts(session.UserId)
		if err != nil {
			fmt.Println(err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		internal.RenderLikePages(w, r, session.Username, likedPosts)
	} else {
		w.WriteHeader(405)
	}
}

func (hh *HttpHandler) HandleDislikedPosts(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		sessionCookie, err := r.Cookie("session_id")
		if err != nil {
			fmt.Println(err)
			http.Redirect(w, r, "/access_denied", http.StatusSeeOther)
			return
		}
		sessionID := sessionCookie.Value

		session, err := hh.business.Session(sessionID)
		if err != nil {
			fmt.Println(err)
			http.Redirect(w, r, "/access_denied", http.StatusSeeOther)
			return
		}

		dislikedPosts, err := hh.business.GetDislikedPosts(session.UserId)
		if err != nil {
			fmt.Println(err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		internal.RenderDislikePages(w, r, session.Username, dislikedPosts)
	} else {
		w.WriteHeader(405)
	}
}
