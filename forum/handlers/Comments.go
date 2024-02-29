package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"forum/forum/domain"
)

func (hh *HttpHandler) HandleAddComment(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		username, err := hh.GetUsername(w, r)
		if err != nil {
			http.Redirect(w, r, "/access_denied", http.StatusSeeOther)
			return
		}

		err = r.ParseForm()
		if err != nil {
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}

		postIDStr := r.Form.Get("post_id")
		commentText := r.Form.Get("comment_text")
		userIdStr := r.Form.Get("user_id")
		postID, err := strconv.Atoi(postIDStr)
		if err != nil {

			http.Error(w, "Invalid post ID", http.StatusBadRequest)
			return
		}
		userId, err := strconv.Atoi(userIdStr)
		if err != nil {

			http.Error(w, "Invalid user ID", http.StatusBadRequest)
			return
		}
		if len(strings.TrimSpace((commentText))) <= 0 {
			fmt.Fprintf(w, "Comments text is required")
			return
		}
		comment := domain.Comments{
			PostId:       postID,
			UserId:       userId,
			Username:     username.Username,
			Content:      commentText,
			CreationDate: time.Now(),
		}

		err = hh.business.AddComment(comment, userId, postID)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/post/?id="+postIDStr, http.StatusSeeOther)
	} else if r.Method == http.MethodDelete {
		postID := r.URL.Query().Get("id")
		if postID == "" {
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}
		id, err := strconv.Atoi(postID)
		if err != nil {
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}

		err = hh.business.DeleteComment(id)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/history", http.StatusSeeOther)

	} else {
		w.WriteHeader(405)
	}
}

func (hh *HttpHandler) DeleteCommentHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		commentID := r.PostFormValue("comment_id")
		if commentID == "" {
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}

		id, err := strconv.Atoi(commentID)
		if err != nil {
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}

		err = hh.business.DeleteComment(id)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/history", http.StatusSeeOther)
		return
	}
}
