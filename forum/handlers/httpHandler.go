package handlers

import (
	"crypto/rand"
	"fmt"
	"net/http"
	"strings"
	"time"

	"forum/forum"
	"forum/forum/domain"

	"golang.org/x/crypto/bcrypt"
)

type HttpHandler struct {
	business forum.Business
}

func NewHandler(Business forum.Business) (HttpHandler, error) {
	return HttpHandler{
		business: Business,
	}, nil
}

func (hh *HttpHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/login":
		hh.HandleUserLogin(w, r)
	case "/register":
		hh.HandleUserRegistration(w, r)

	case "/":
		hh.MainHandler(w, r)
	case "/static/style.css":
		hh.CssHandler(w, r)
	case "/createPost":
		hh.HandleNewPost(w, r)
	case "/add_comment":
		hh.HandleAddComment(w, r)
	case "/my_posts":
		hh.HandleMyPosts(w, r)
	case "/notifications":
		hh.NotificationHandler(w, r)
	case "/liked_posts":
		hh.HandleLikedPosts(w, r)
	case "/history":
		hh.UserActivityHandler(w, r)
	case "/filtered-posts":
		hh.HandleFilteredPosts(w, r)
	case "/like_dislike_post":
		hh.HandleLikeDislikePost(w, r)
	case "/like_dislike_comment":
		hh.HandleLikeDislikeComment(w, r)
	case "/disliked_posts":
		hh.HandleDislikedPosts(w, r)
	case "/delete_comment":
		hh.DeleteCommentHandler(w, r)
	case "/exit":
		hh.HandleLogout(w, r)
	case "/access_denied":
		hh.Handle403(w, r)
	default:
		if strings.HasPrefix(r.URL.Path, "/post/") {
			hh.HandlePostDetails(w, r)
			return
		}
		if strings.HasPrefix(r.URL.Path, "/edit_post/") {
			hh.HandleEditPostDetails(w, r)

			return
		}
		if strings.HasPrefix(r.URL.Path, "/comment/") {
			hh.HandleCommentEdit(w, r)

			return
		}
		hh.Handle404(w, r)
	}
}

func (hh *HttpHandler) CssHandler(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "text/css")
	http.ServeFile(writer, request, "./forum/static/style.css")
}

func generateUniqueFilename() string {
	randomString := generateRandomString(10)

	timestamp := time.Now().UnixNano()

	combinedString := fmt.Sprintf("%s%d", randomString, timestamp)

	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(combinedString), bcrypt.DefaultCost)
	if err != nil {
		return ""
	}

	hashHex := fmt.Sprintf("%x", hashedBytes)

	return hashHex
}

func generateRandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	randomBytes := make([]byte, length)
	_, err := rand.Read(randomBytes)
	if err != nil {
		return ""
	}

	for i, b := range randomBytes {
		randomBytes[i] = charset[b%byte(len(charset))]
	}

	return string(randomBytes)
}

func (hh *HttpHandler) GetUsername(w http.ResponseWriter, r *http.Request) (*domain.Session, error) {
	sessionCookie, err := r.Cookie("session_id")
	if err != nil {
		return nil, err
	}
	sessionID := sessionCookie.Value

	session, err := hh.business.Session(sessionID)
	if err != nil {
		return nil, err
	}
	return session, nil
}

func formatTimestamp(timestamp time.Time) string {
	T := timestamp.Format("2006-01-02 15:04:05")

	return T
}
