package main

import (
	"cloud.google.com/go/firestore"
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
)

const version = "1.0.0"

// Define an application struct to hold the dependencies for our HTTP handlers, helpers,
// and middleware.
type application struct {
	logger     *log.Logger
	dbConfig   DBConfig
	dbSnippets SnippetModeler
}

type DBConfig struct {
	projectID  string
	databaseID string
	collection string
}

func main() {

	projectID := os.Getenv("GOOGLE_CLOUD_PROJECT")
	if projectID == "" {
		log.Fatalf("Environment variable not set %v", projectID)
	}

	databaseID := os.Getenv("DATABASE_ID")
	if databaseID == "" {
		log.Fatalf("Environment variable not set %v", databaseID)
	}

	collection := os.Getenv("COLLECTION_NAME")
	if collection == "" {
		log.Fatalf("Environment variable not set %v", collection)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("defaulting to port %s", port)
	}

	dbConfig := DBConfig{
		projectID:  projectID,
		databaseID: databaseID,
		collection: collection,
	}

	// Initialize a new logger which writes messages to the standard out stream,
	// prefixed with the current date and time.
	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	// Declare an instance of the application struct, containing the config struct and
	// the logger.
	app := &application{
		logger:   logger,
		dbConfig: dbConfig,
	}

	db, err := app.openDB()

	if err != nil {
		logger.Fatal(err)
	}

	defer db.Close()

	app.dbSnippets = &SnippetModel{
		Client:     db,
		Collection: collection,
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/home", app.home)
	mux.HandleFunc("/v1/healthcheck", app.healthcheckHandler)
	mux.HandleFunc("/find-edm/all-events", app.GetAllEvents)
	mux.HandleFunc("/find-edm/venue", app.GetEdmEventByClubName)
	mux.HandleFunc("/find-edm/artist", app.GetEdmEventByArtistName)

	// Start the HTTP server.
	logger.Printf("starting server on %s", port)
	err = http.ListenAndServe(":"+port, mux)
	if err != nil {
		logger.Fatal(err)
	}
}

func (app *application) openDB() (*firestore.Client, error) {

	ctx := context.Background()

	client, err := firestore.NewClientWithDatabase(ctx, app.dbConfig.projectID, app.dbConfig.databaseID)

	if err != nil {
		return nil, err
	}

	fmt.Println("Connected to Firestore!")

	return client, nil
}
