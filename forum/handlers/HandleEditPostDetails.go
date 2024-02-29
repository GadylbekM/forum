package handlers

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"forum/forum/domain"
	"forum/forum/internal"
)

func (hh *HttpHandler) HandleEditPostDetails(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		postIDStr := r.URL.Query().Get("id")
		postID, err := strconv.Atoi(postIDStr)
		if err != nil {
			fmt.Println(err)
			hh.Handle404(w, r)
			return
		}

		sessionCookie, err := r.Cookie("session_id")
		if err != nil {
			hh.Handle403(w, r)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)

			return
		}
		sessionID := sessionCookie.Value

		session, err := hh.business.Session(sessionID)
		if err != nil || session == nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		category := r.FormValue("category")
		title := r.FormValue("title")
		content := r.FormValue("content")

		r.ParseMultipartForm(20 << 20)
		file, _, err := r.FormFile("image")
		if err != nil && err != http.ErrMissingFile {
			http.Error(w, "Error uploading image", http.StatusBadRequest)
			return
		}

		if file != nil {
			defer file.Close()
		}
		var imagePathHTML string

		if file == nil {
			imagePathHTML = "../static/Gallery/default.png"
		} else {

			var fileSize int64
			buf := make([]byte, 1024)

			for {
				n, err := file.Read(buf)
				if err != nil && err != io.EOF {
					http.Error(w, "Error reading file", http.StatusInternalServerError)
					return
				}
				if n == 0 {
					break
				}
				fileSize += int64(n)
			}

			const maxFileSize = 20 << 20
			if fileSize > maxFileSize {
				internal.RenderEditPostPage(w, r, session.Username, "Image file size exceeds the limit (20MB) ", postIDStr)
				return
			}

			_, err := file.Seek(0, io.SeekStart)
			if err != nil {
				http.Error(w, "Error reading file", http.StatusInternalServerError)
				return
			}

			imageFilename := generateUniqueFilename()
			imagePath := "./forum/static/Gallery/" + imageFilename + ".png"

			f, err := os.Create(imagePath)
			if err != nil {
				http.Error(w, "Error saving image", http.StatusInternalServerError)
				return
			}
			defer f.Close()

			imagePathHTML = "../static/Gallery/" + imageFilename + ".png"

			_, err = io.Copy(f, file)
			if err != nil {
				http.Error(w, "Error saving image", http.StatusInternalServerError)
				return
			}
		}
		if len(strings.TrimSpace(category)) <= 0 || len(strings.TrimSpace(title)) <= 0 {
			internal.RenderEditPostPage(w, r, session.Username, "The post title and content must not be empty", postIDStr)
			return
		}
		newPost := domain.Posts{
			Username:     session.Username,
			UserId:       session.UserId,
			Category:     category,
			Title:        title,
			Content:      content,
			CategoryId:   1,
			ImageField:   imagePathHTML,
			CreationDate: time.Now(),
		}

		err = hh.business.EditPost(postID, newPost)
		if err != nil {
			fmt.Print(err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/", http.StatusSeeOther)

	} else if r.Method == http.MethodGet {
		postIDStr := r.URL.Query().Get("id")
		username, err := hh.GetUsername(w, r)
		if err != nil {
			hh.Handle404(w, r)
			return
		}
		posts, err := hh.business.GetAllPosts()
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		if username.Username == "" {

			internal.RenderMainPage(w, r, username, posts)
			return
		}

		err = internal.RenderEditPostPage(w, r, username.Username, "", postIDStr)
		if err != nil {
			fmt.Println(err)
		}
	} else {
		w.WriteHeader(405)
	}
}
