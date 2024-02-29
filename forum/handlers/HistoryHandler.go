package handlers

import (
	"fmt"
	"net/http"

	"forum/forum/internal"
)

func (hh *HttpHandler) UserActivityHandler(w http.ResponseWriter, r *http.Request) {
	username, err := hh.GetUsername(w, r)
	if err != nil {
		// Handle the error, e.g., by displaying an error page

		http.Redirect(w, r, "/access_denied", http.StatusSeeOther)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if username.Username == "" {
		// Handle the case when the user is not logged in
		http.Redirect(w, r, "/access_denied", http.StatusSeeOther)
		return
	}

	// Retrieve the user's created posts and comments
	userActivity, err := hh.business.GetUserActivity(username.UserId)
	if err != nil {
		// Handle the error, e.g., by displaying an error page
		fmt.Println(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Render the user activity page
	internal.RenderUserActivityPage(w, r, username.Username, userActivity)
}
