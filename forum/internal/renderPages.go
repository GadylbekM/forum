package internal

import (
	"fmt"
	"html/template"
	"net/http"

	"forum/forum/domain"
)

func RenderMainPage(w http.ResponseWriter, r *http.Request, userSession *domain.Session, posts []domain.Posts) {
	tmpl, err := template.ParseFiles("./forum/templates/index.html", "./forum/templates/base.html")
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	if userSession != nil {
	}
	data := struct {
		Name  string
		Posts []domain.Posts
	}{
		Posts: posts,
	}
	if userSession == nil {
		data.Name = "Guest"
	} else {
		data.Name = userSession.Username
	}

	err = tmpl.Execute(w, data)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func RenderUserActivityPage(w http.ResponseWriter, r *http.Request, username string, activity domain.UserActivity) {
	tmpl, err := template.ParseFiles("./forum/templates/History.html", "./forum/templates/base.html")
	if err != nil {
		fmt.Println("Cant get the HTML files")

		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Execute the template and write the response
	data := struct {
		Name         string
		CreatedPosts []domain.Posts
		Comments     []domain.Comments
	}{
		Name:         username,
		CreatedPosts: activity.CreatedPosts,
		Comments:     activity.Comments,
	}

	err = tmpl.Execute(w, data)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func RenderNotifications(w http.ResponseWriter, username string, notifications []domain.Notification, notifications_comment []domain.Notification_comments) {
	tmpl, err := template.ParseFiles("./forum/templates/notify.html", "./forum/templates/base.html")
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	data := struct {
		Name                  string
		Notifications         []domain.Notification
		Notification_comments []domain.Notification_comments
	}{
		Name:                  username,
		Notifications:         notifications,
		Notification_comments: notifications_comment,
	}

	err = tmpl.Execute(w, data)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func RenderLoginPage(w http.ResponseWriter, r *http.Request, errorMessage string) {
	tmpl, err := template.ParseFiles("./forum/templates/login.html")
	if err != nil {
		fmt.Println("Cant get the HTML files")
		return
	}
	data := struct {
		ErrorMessage string
	}{
		ErrorMessage: errorMessage,
	}

	err = tmpl.Execute(w, data)

	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func RenderRegisterPage(w http.ResponseWriter, r *http.Request, errorMessage string) {
	tmpl, err := template.ParseFiles("./forum/templates/register.html")
	if err != nil {
		fmt.Println("Cant get the HTML files")
		return
	}
	data := struct {
		ErrorMessage string
	}{
		ErrorMessage: errorMessage,
	}

	err = tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func RenderLikePages(w http.ResponseWriter, r *http.Request, username string, posts []domain.Posts) {
	tmpl, err := template.ParseFiles("./forum/templates/likedPosts.html", "./forum/templates/base.html")
	if err != nil {
		fmt.Println(err)
		fmt.Println("dada")

		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	data := struct {
		Username string
		Posts    []domain.Posts
	}{
		Username: username,
		Posts:    posts,
	}
	err = tmpl.Execute(w, data)
	if err != nil {
		fmt.Println(err)
		fmt.Println("dada")

		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func RenderDislikePages(w http.ResponseWriter, r *http.Request, username string, dislikedposts []domain.Posts) {
	tmpl, err := template.ParseFiles("./forum/templates/DislikedPosts.html", "./forum/templates/base.html")
	if err != nil {
		fmt.Println(err)
		fmt.Println("dada")

		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	data := struct {
		Username     string
		DislikedPost []domain.Posts
	}{
		Username:     username,
		DislikedPost: dislikedposts,
	}
	err = tmpl.Execute(w, data)
	if err != nil {
		fmt.Println(err)
		fmt.Println("dada")

		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func RenderPostPage(w http.ResponseWriter, r *http.Request, username string, error string) {
	tmpl, err := template.ParseFiles("./forum/templates/createPost.html", "./forum/templates/base.html")
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	data := struct {
		Username string
		Error    string
	}{
		Username: username,
		Error:    error,
	}

	err = tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func RenderEditPostPage(w http.ResponseWriter, r *http.Request, username string, error string, postID string) error {
	tmpl, err := template.ParseFiles("./forum/templates/Edit_Post.html", "./forum/templates/base.html")
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return err
	}

	data := struct {
		Username string
		Error    string
		PostId   string
	}{
		Username: username,
		Error:    error,
		PostId:   postID,
	}
	err = tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return err
	}
	return nil
}

func RenderEditCommentPage(w http.ResponseWriter, r *http.Request, username string, error string, commentID string) error {
	tmpl, err := template.ParseFiles("./forum/templates/Edit_Comment.html", "./forum/templates/base.html")
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return err
	}

	data := struct {
		Username  string
		Error     string
		CommentId string
	}{
		Username:  username,
		Error:     error,
		CommentId: commentID,
	}
	err = tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return err
	}
	return nil
}

func RenderMyPostPage(w http.ResponseWriter, r *http.Request, username string, posts []domain.Posts) {
	tmpl, err := template.ParseFiles("./forum/templates/my_posts.html", "./forum/templates/base.html")
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	data := struct {
		Username string
		Posts    []domain.Posts
	}{
		Username: username,
		Posts:    posts,
	}
	err = tmpl.Execute(w, data)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func RenderErrorPage(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("./forum/templates/404.html")
	if err != nil {
		fmt.Println("Cant get the HTML files")
		return
	}

	err = tmpl.Execute(w, nil)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func RenderError403Page(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("./forum/templates/403.html")
	if err != nil {
		fmt.Println("Cant get the HTML files")
		return
	}

	err = tmpl.Execute(w, nil)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func RenderAboutPage(w http.ResponseWriter, r *http.Request, userSession *domain.Session, posts domain.Posts, comments []domain.Comments) {
	tmpl, err := template.ParseFiles("./forum/templates/About.html", "./forum/templates/base.html")
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	data := struct {
		Name     string
		UserId   int
		Post     domain.Posts
		Comments []domain.Comments
	}{
		Post:     posts,
		Comments: comments,
	}
	if userSession == nil {
		data.Name = "Guest"
		data.UserId = 0
	} else {
		data.Name = userSession.Username
		data.UserId = userSession.UserId
	}

	err = tmpl.Execute(w, data)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}
