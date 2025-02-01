package server

import (
	"log"
	"net/http"
	"strconv"

	strct "reboot01.com/js/forum/structs"
)

func NotificationsPage(w http.ResponseWriter, r *http.Request) {
	userID, err := strconv.Atoi(r.FormValue("user"))
	if err != nil {
		log.Println("Error converting userID to int:", err)
		err := strct.ErrorPageData{Code: "500", ErrorMsg: "INTERNAL SERVER ERROR"}
		errHandler(w, r, &err)
		return
	}

	data := struct {
		UserID        int
		Avatar        string
		Notifications []strct.Notification
	}{
		UserID: userID,
		// Avatar:        session.Values["avatar"].(string),
		// Notifications: notifications,
	}

	err = strct.Templates.ExecuteTemplate(w, "notifications.html", data)
	if err != nil {
		log.Println("Error rendering notifications page:", err)
		err := strct.ErrorPageData{Code: "500", ErrorMsg: "INTERNAL SERVER ERROR"}
		errHandler(w, r, &err)
		return
	}
}
