package api

import (
	"net/http"

	"reboot01.com/js/realtime-forum/internal/response"
)

func (app *Application) fetchCategories(w http.ResponseWriter, r *http.Request) {
	categories, err := app.DB.GetAllCategories()
	if err != nil {
		app.serverError(w, r, err)
		return
	}
	err = response.JSON(w, http.StatusOK, categories)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

}
