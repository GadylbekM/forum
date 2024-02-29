package handlers

import (
	"fmt"
	"net/http"
	"strconv"
)

func (hh *HttpHandler) HandleLikeDislikePost(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	session, err := hh.GetUsername(w, r)
	if err != nil {
		http.Redirect(w, r, "/access_denied", http.StatusSeeOther)
		return
	}

	postIDStr := r.FormValue("post_id")

	ownerID, err := strconv.Atoi(r.FormValue("owner_id"))
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	postID, err := strconv.Atoi(postIDStr)
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	action := r.FormValue("action")

	switch action {
	case "like":
		err = hh.business.LikePost(postID, session.UserId, ownerID, "like your post", session.Username)
	case "dislike":
		err = hh.business.DislikePost(postID, session.UserId, ownerID, "dislike your post", session.Username)
	default:
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	if err != nil {
		hh.Handle404(w, r)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	redirectURL := r.Referer()
	http.Redirect(w, r, redirectURL, http.StatusSeeOther)
}

func (hh *HttpHandler) HandleLikeDislikeComment(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	session, err := hh.GetUsername(w, r)
	if err != nil {
		http.Redirect(w, r, "/access_denied", http.StatusSeeOther)
		return
	}

	ownerID, err := strconv.Atoi(r.FormValue("owner_id"))
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	postID, err := strconv.Atoi(r.FormValue("post_id"))
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	commentIDStr := r.FormValue("comment_id")
	commentID, err := strconv.Atoi(commentIDStr)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	action := r.FormValue("action")

	switch action {
	case "like":
		err = hh.business.LikeComment(commentID, session.UserId, ownerID, "likes your comment", session.Username, postID)
		if err != nil {
			fmt.Println(err)
		}
	case "dislike":
		err = hh.business.DislikeComment(commentID, session.UserId, ownerID, "dislikes your comment", session.Username, postID)
	default:
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	if err != nil {
		hh.Handle404(w, r)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	redirectURL := r.Referer()
	http.Redirect(w, r, redirectURL, http.StatusSeeOther)
}
