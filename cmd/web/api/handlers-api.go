package main

import (
	"awesomeProject1/internal/cards"
	"encoding/json"
	"net/http"
	"strconv"
)

// Used in the GetPaymentIntent function below.
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
	// Declaring the payload variable as type stripePayload. (from the struct above)
	var payload stripePayload

	// Retrieve the body of the request and decode it into the variable payload.
	// The payload variable will be populated with the currency and amount.
	// err is assigned the value of: from the JSON pkg, get a new decoder,
	// and decode the request body into a reference to payload.
	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		app.errorLog.Println(err)
		return
	}

	// The currency is a string and the amount is a string.
	// Need them to be a numeric value. So use Atoi: Alpha to Numeric.
	// Decode body of that request into that variable payload
	amount, err := strconv.Atoi(payload.Amount)
	if err != nil {
		app.errorLog.Println(err)
		return
	}

	/*
		Utilize the cards dot card type:
		Create a variable, named card, and use the cards package to assign it's value.
		And populate it's fields. Gain access to needed methods
		The values for Secret and Key will be pulled from app.config.
	*/
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
		if err != nil {                               // Handle an error.
			app.errorLog.Println(err)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(out)

		// But if things don't work as expected.
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
