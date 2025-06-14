package api

import (
	"fmt"
	"net/http"
)

func (app *Application) backgroundTask(r *http.Request, fn func() error) {
	app.WG.Add(1)

	go func() {
		defer app.WG.Done()

		defer func() {
			err := recover()
			if err != nil {
				app.reportServerError(r, fmt.Errorf("%s", err))
			}
		}()

		err := fn()
		if err != nil {
			app.reportServerError(r, err)
		}
	}()
}
