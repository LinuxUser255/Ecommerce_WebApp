package main

import (
	"encoding/json"
	"myapp/internal/cards"
	"net/http"
	"strconv"
)

type stripePayload struct {
	Currency string `json:"currency"`
	Amount   string `json:"amount"`
}

type jsonResponse struct {
	OK      bool   `json:"ok"`
	Message string `json:"message,omitempty"`
	Content string `json:"content,omitempty"`
	ID      int    `json:"id,omitempty"`
}

// GetPaymentIntent This request will be called using POST, and contain the request body in json format.
// That json file, will conform to the standard set in stripePayload.
func (app *application) GetPaymentIntent(w http.ResponseWriter, r *http.Request) {
	var payload stripePayload

	// Decode body of that request into that variable payload.
	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		app.errorLog.Println(err)
		return
	}

	// Convert the payload data to Alphanumeric format.
	amount, err := strconv.Atoi(payload.Amount)
	if err != nil {
		app.errorLog.Println(err)
		return
	}

	// Gain access to needed methods
	// The values for Secret and Key will be pulled from app.config.
	card := cards.Card{
		Secret:   app.config.stripe.secret,
		Key:      app.config.stripe.key,
		Currency: payload.Currency,
	}

	/*
	   Declare a variable, name it OK, with the assumption that it will be true
	   Then, call the function charge in cards and get the three variables that are returns.

	   Payment Intent (pi)
	   message (msg)
	   error (err)

	   Populate them by calling card.Charge() and requires 2 vars the currency and the amount.
	*/
	okay := true
	pi, msg, err := card.Charge(payload.Currency, amount)
	if err != nil {
		okay = false
	}

	if okay {
		out, err := json.MarshalIndent(pi, "", "   ") // Marshal returns the JSON encoding.
		if err != nil {                               
			app.errorLog.Println(err)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(out)

		// But if things don't work..
	} else {
		j := jsonResponse{
			OK:      false,
			Message: msg,
			Content: "",
		}

		out, err := json.MarshalIndent(j, "", "   ")
		if err != nil {
			app.errorLog.Println(err)
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(out)
	}
}
