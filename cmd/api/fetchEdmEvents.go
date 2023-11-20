package main

import (
	"fmt"
	"net/http"
)

const QUERY_CLUBNAME = "clubname"
const QUERY_ARTISTNAME = "artistname"
const QUERY_EVENTDATE = "eventdate"

func (app *application) GetAllEvents(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	edmEvents, err := app.dbSnippets.FindAll()

	if err != nil {
		http.Error(w, "The server encountered a problem and could not process your request", http.StatusInternalServerError)
	}

	if edmEvents != nil {
		err = app.writeJSON(w, http.StatusOK, edmEvents, nil)
	} else {
		e := EdmEventNotFound{Description: "Search Criteria did not return anything"}
		err = app.writeJSON(w, http.StatusNotFound, e, nil)
	}

	if err != nil {
		app.logger.Print(err)
		http.Error(w, "The server encountered a problem and could not process your request", http.StatusInternalServerError)
	}

}

func (app *application) GetEdmEventByClubName(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	query := r.URL.Query().Get(QUERY_CLUBNAME)
	fmt.Println(query)
	edmEvents, err := app.dbSnippets.FindAllWithFilter(QUERY_CLUBNAME, query)

	if err != nil {
		http.Error(w, "The server encountered a problem and could not process your request", http.StatusInternalServerError)
	}

	if edmEvents != nil {
		err = app.writeJSON(w, http.StatusOK, edmEvents, nil)
	} else {
		e := EdmEventNotFound{Description: "Search Criteria did not return anything"}
		err = app.writeJSON(w, http.StatusNotFound, e, nil)
	}

	if err != nil {
		app.logger.Print(err)
		http.Error(w, "The server encountered a problem and could not process your request", http.StatusInternalServerError)
	}

}

func (app *application) GetEdmEventByArtistName(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	query := r.URL.Query().Get(QUERY_ARTISTNAME)
	fmt.Println(query)
	edmEvents, err := app.dbSnippets.FindAllWithFilter(QUERY_ARTISTNAME, query)

	if err != nil {
		http.Error(w, "The server encountered a problem and could not process your request", http.StatusInternalServerError)
	}

	if edmEvents != nil {
		err = app.writeJSON(w, http.StatusOK, edmEvents, nil)
	} else {
		e := EdmEventNotFound{Description: "Search Criteria did not return anything"}
		err = app.writeJSON(w, http.StatusNotFound, e, nil)
	}

	if err != nil {
		app.logger.Print(err)
		http.Error(w, "The server encountered a problem and could not process your request", http.StatusInternalServerError)
	}

}

/// just need to add the const QUERY_EVENTDATE and this portion should be wrapped up

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
}
