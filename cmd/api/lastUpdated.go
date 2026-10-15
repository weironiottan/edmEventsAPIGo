package main

import (
	"net/http"
	"time"
)

type LastUpdatedResponse struct {
	LastUpdated string `json:"lastUpdated"`
}

// GetLastUpdated serves when the scraper last refreshed the events, e.g.
// {"lastUpdated":"2026-10-14T09:00:00Z"}
func (app *application) GetLastUpdated(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	lastUpdated, found, err := app.dbSnippets.LastUpdated()

	if err != nil {
		app.logger.Print(err)
		http.Error(w, "The server encountered a problem and could not process your request", http.StatusInternalServerError)
		return
	}

	if !found {
		e := EdmEventNotFound{Description: "No scrape has been recorded yet"}
		err = app.writeJSON(w, http.StatusNotFound, e, nil)
	} else {
		// The scraper runs daily, so a few minutes of caching is plenty fresh
		headers := http.Header{"Cache-Control": []string{"public, max-age=300"}}
		response := LastUpdatedResponse{LastUpdated: lastUpdated.UTC().Format(time.RFC3339)}
		err = app.writeJSON(w, http.StatusOK, response, headers)
	}

	if err != nil {
		app.logger.Print(err)
		http.Error(w, "The server encountered a problem and could not process your request", http.StatusInternalServerError)
	}
}
