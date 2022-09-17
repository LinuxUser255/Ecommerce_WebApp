package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

// The css version used for appending any external CSS or Java Script files.
// This will force most browsers to revert to a lower version if incremented.
const VERSION = "1.0.0"

// Creating the config type: Holds configuration info for the app:
// The port number, the api  string: url used to make backend api calls
// dsn string is database, data source name: how to connect to the DB
type Config struct {
	port int
	env  string
	db   struct {
		dsn string
	}
	stripe struct {
		secret string
		key    string
	}
}

// Creating the receiver type.
type application struct {
	config   Config
	infoLog  *log.Logger
	errorLog *log.Logger
	version  string
}

// Create the Web Server.
// Has the receiver of "app" of type pointer to Application
// Assigning a value to the Server variable, and calling it from the HTTP package: http.server
func (app *application) serve() error {
	srv := &http.Server{
		Addr:              fmt.Sprintf(":%d", app.config.port),
		Handler:           app.routes(),
		IdleTimeout:       30 * time.Second,
		ReadTimeout:       10 * time.Second,
		ReadHeaderTimeout: 5 * time.Second,
		WriteTimeout:      5 * time.Second,
	}

	// Prints a message to the terminal telling the user that the Server has started.
	app.infoLog.Println(fmt.Sprintf("Starting Back End server in %s mode on port %d", app.config.env, app.config.port))

	// Returns an error message if something went wrong. And nothing if all worked correctly.
	return srv.ListenAndServe()
}

func main() {
	var cfg Config

	flag.IntVar(&cfg.port, "port", 4001, "Server port to listen on")
	flag.StringVar(&cfg.env, "env", "development", "Application environment {development|production|maintenance}")

	// Parse the commands.
	flag.Parse()

	/* ATTN! Security issue addressed and securely coded.
		Read the following commented lines for details.
		Defining the Stripe key and secret configuration.
	 	Get the Stripe publishable and private key, but don't make them visible on the
	 	command line.
		This would introduce an Information Disclosure vulnerability.
	 	The private key could be revealed by executing the command: ps -ax or ps -aux
	 	To prevent this, read the Stripe key and Stripe Secrete, from the environment
	 	variables */
	cfg.stripe.key = os.Getenv("STRIPE_KEY")       // <-- This gets the stripe key from environment
	cfg.stripe.secret = os.Getenv("STRIPE_SECRET") // <-- This gets the stripe secret from environment.

	// Creating a log to display information and errors.
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	// The application variable.
	app := &application{
		config:   cfg,
		infoLog:  infoLog,
		errorLog: errorLog,
		//	templateCache: make(map[string]*template.Template),
		version: VERSION,
	}

	err := app.serve()
	if err != nil {
		log.Fatal(err)

	}

}
