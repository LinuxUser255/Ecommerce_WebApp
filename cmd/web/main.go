package main

import (
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"
)

// main.go is the primary web server for this app

/*
The css version used for appending any external CSS or Java Script files.
This will force most browsers to revert to a lower version if incremented.
saves trouble of clearing cache
*/
const VERSION = "1.0.0"
const cssVersion = "1"

/*
Creating the config type:
Holds configuration info for the app:
The port number, the api  string, the api url used to make backend api calls
dsn string is database, data source name how to connect to the DB
*/
type Config struct {
	port int
	env  string
	api  string
	db   struct {
		dsn string
	}
	// Your Stripe API Keys
	stripe struct {
		secret string
		key    string
	}
}

/*
Creating the app type:
Holds all the information needed to run the app:
The config type, the http.Server type, and the template type
*/
type App struct {
	config   Config
	server   *http.Server
	template *template.Template
}

/*
Creating the receiver type which includes loggers.
info and error Log are a pointer to log.Logger
Template cache, is a map of type string,
and the content of each entry will be a pointer to template
*/
type application struct {
	config        Config
	infoLog       *log.Logger
	errorLog      *log.Logger
	templateCache map[string]*template.Template
	version       string
}

/*
This "calls" the server: Create the Web Server.
Has the receiver of "app" of type pointer to Application
Assigning a value to the Server variable, and calling it from the HTTP package: http.server
Setting a time-out for the server.
*/
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
	app.infoLog.Println(fmt.Sprintf("Starting HTTP server in %s mode on port %d", app.config.env, app.config.port))

	// Returns an error message if something went wrong. And nothing if all worked correctly.
	return srv.ListenAndServe()
}

/*
The Main function:
Populating the variable with command line flags and arguments.
Read the command flag into the config variable
*/
func main() {
	var cfg Config

	flag.IntVar(&cfg.port, "port", 4000, "Server port  on which to listen")
	flag.StringVar(&cfg.env, "env", "development", "Application enviornment {development|production}")
	flag.StringVar(&cfg.api, "api", "http://localhost:4001", "URL to api")

	// Parse the commands.
	flag.Parse()

	/* ATTN! Security issue addressed and securely coded.
	The private key could be revealed by executing the command: ps -ax or ps -aux
	To prevent this, READ the Stripe key and Stripe Secret, FROM the ENVIRONMENT
	Variables
	*/
	cfg.stripe.key = os.Getenv("STRIPE_KEY")    // <-- Retrieves the stripe key from environment
	cfg.stripe.secret = os.Getenv("STRIPE_KEY") // <--  Retrieves the stripe secret from environment.

	// Creating a log to display information and errors.
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	// Creating a map for the template cache. *pointer to template
	// tc = Template Cache
	tc := make(map[string]*template.Template)

	// Creating an application variable referencing an application.
	// Listed fields refer to variables created above.
	app := &application{
		config:        cfg,
		infoLog:       infoLog,
		errorLog:      errorLog,
		templateCache: tc,
		version:       VERSION,
	}

	// Error message to be created if there is a problem running the application server
	err := app.serve()
	if err != nil {
		app.errorLog.Println(err)
		log.Fatal(err)
	}

}
