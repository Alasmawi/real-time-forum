package server

import (
	strct "reboot01.com/js/forum/structs"
	"log"
	"net/http"
)

func errHandler(w http.ResponseWriter, _ *http.Request, errData *strct.ErrorPageData) {
	err := strct.Templates.ExecuteTemplate(w, "error.html", errData)
	if err != nil {
		log.Println("Error rendering error page:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	// errCodeInt, _ := strconv.Atoi(errData.Code)

	// w.WriteHeader(errCodeInt)
}

func AutherrHandler(w http.ResponseWriter, _ *http.Request, errData *strct.ErrorPageData) {
	err := strct.Templates.ExecuteTemplate(w, "error.html", errData)
	if err != nil {
		log.Println("Error rendering error page:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
