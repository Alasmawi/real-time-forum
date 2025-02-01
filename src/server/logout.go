package server

import (
	"log"
	"net/http"
	"time"

	strct "reboot01.com/js/forum/structs"
)

func Logout(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("userID")

	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    "",
		Expires:  time.Now().Add(-time.Hour),
		HttpOnly: true,
	})

	_, err := strct.Db.Exec("DELETE FROM session WHERE user_id = ?", userID)
	if err != nil {
		log.Fatal(err)
	}

	_, err = strct.Db.Exec("UPDATE user SET current_session = NULL WHERE user_id = ?", userID)
	if err != nil {
		log.Fatal(err)
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
