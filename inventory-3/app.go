package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

type App struct {
	DB		*sql.DB
	Router	*mux.Router
}

func (app *App) Initialise() error {
	connnectString := fmt.Sprintf("%v:%v@tcp(127.0.0.1:3306)/%v", DbUser, DbPassword, DbName)
	var err error
	app.DB, err = sql.Open("mysql", connnectString)
	if err != nil {
		return err
	}

	err = app.DB.Ping()
	if err != nil {
		app.DB.Close()
		return err
	}

	app.Router = mux.NewRouter().StrictSlash(true)
	app.handleRouters()
	return nil

}

func (app *App) Run(address string) {
	log.Fatal(http.ListenAndServe(address, app.Router))
}


func (app *App) handleRouters() {
	app.Router.HandleFunc("/homepage", app.homepage)
	app.Router.HandleFunc("/movies", app.getMovies).Methods("GET")
	app.Router.HandleFunc("/movies/{id}", app.getMovie).Methods("GET")

}

func sendError(w http.ResponseWriter, statusCode int, err string) {
	error_message := map[string]string{"error": err}
	sendResponse(w, statusCode, error_message)
}

func sendResponse(w http.ResponseWriter, statusCode int, payLoad interface{}) {
	response, err := json.Marshal(payLoad)
	if err != nil {
		sendError(w, http.StatusInternalServerError, err.Error())
	}
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(statusCode)
	w.Write(response)
}

func (app *App) homepage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "text/plain")
	fmt.Fprintln(w, "welcome to movies Home Page!")
	log.Println("endpint hit: homepage")
}

func (app *App) getMovies(w http.ResponseWriter, r *http.Request) {
	movies, err := getMovies(app.DB)
	if err != nil {
		sendError(w, http.StatusInternalServerError, err.Error())
	}
	sendResponse(w, http.StatusOK, movies)
}

func (app *App) getMovie(w http.ResponseWriter, r * http.Request) {
	vars := mux.Vars(r)
	key, err := strconv.Atoi(vars["id"])
	if err != nil {
		sendError(w, http.StatusBadRequest, "invalid movie id")
		return
	}
	m := Movie{Id: key}
	err = m.getMovies(app.DB)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			sendError(w, http.StatusNotFound, "movie not found")
		default:
			sendError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}
	sendResponse(w, http.StatusOK, m)
}