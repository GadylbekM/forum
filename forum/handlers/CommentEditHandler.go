package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"forum/forum/domain"
	"forum/forum/internal"
)

func (hh *HttpHandler) HandleCommentEdit(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		commentIDStr := r.URL.Query().Get("id")
		commentID, err := strconv.Atoi(commentIDStr)
		if err != nil {
			fmt.Println(err)
			hh.Handle404(w, r)
			return
		}

		sessionCookie, err := r.Cookie("session_id")
		if err != nil {
			hh.Handle404(w, r)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)

			return
		}
		sessionID := sessionCookie.Value

		session, err := hh.business.Session(sessionID)
		if err != nil || session == nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		commentText := r.FormValue("comment_text")
		newComment := domain.Comments{
			CommentId: commentID,
			Content:   commentText,
		}

		err = hh.business.EditComment(commentID, newComment)
		if err != nil {
			fmt.Print(err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/", http.StatusSeeOther)

	} else if r.Method == http.MethodGet {
		commentIDStr := r.URL.Query().Get("id")
		username, err := hh.GetUsername(w, r)
		if err != nil {
			hh.Handle404(w, r)
			return
		}
		posts, err := hh.business.GetAllPosts()
		if err != nil {
			fmt.Println("Cant get Posts")
			return
		}
		if username.Username == "" {

			internal.RenderMainPage(w, r, username, posts)
			return
		}

		err = internal.RenderEditCommentPage(w, r, username.Username, "", commentIDStr)
		if err != nil {
			fmt.Println(err)
		}

	} else {
		w.WriteHeader(405)
	}
}
