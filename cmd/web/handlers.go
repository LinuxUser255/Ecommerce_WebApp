package main

import (
	"net/http"
)

func (app *application) PaymentPortal(w http.ResponseWriter, r *http.Request) {
	if err := app.renderTemplate(w, r, "portal", nil); err != nil {
		app.errorLog.Println(err)
	}
}
