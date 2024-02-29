package handlers

import (
	"fmt"
	"net/http"

	"forum/forum/internal"
)

func (hh *HttpHandler) NotificationHandler(w http.ResponseWriter, r *http.Request) {
	username, err := hh.GetUsername(w, r)
	if err != nil {
		// Handle the error, e.g., by displaying an error page

		http.Redirect(w, r, "/access_denied", http.StatusSeeOther)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	notifications, err := hh.business.GetAllNotifications(username.UserId)
	if err != nil {
		fmt.Println(err)
		return
	}
	notifications_comments, err := hh.business.GetAllNotificationsComment(username.UserId)
	if err != nil {
		fmt.Println(err)
		return
	}

	internal.RenderNotifications(w, username.Username, notifications, notifications_comments) // Replace "Username" with the actual username
}
