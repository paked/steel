package main

import (
	"fmt"
	"net/http"

	// "github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
	"github.com/paked/gerrycode/communicator"
	"github.com/paked/steel/models"
)

func main() {
	routes := mux.NewRouter()

	m := routes.PathPrefix("/api").Subrouter()
	m.HandleFunc("/user/login", LoginHandler).Methods("POST")
	m.HandleFunc("/user/register", RegisterHandler).Methods("POST")

	r := m.PathPrefix("/").Subrouter()
	r.HandleFunc("/party", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "you got into the party!")
	})

	models.InitDB()

	http.ListenAndServe("localhost:8080", routes)
	/*
		n := negroni.New()
		n.Use(negroni.NewLogger())
		n.Use(negroni.NewStatic(http.Dir("static/")))


		n.Run("localhost:8080")*/
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	c := communicator.New(w)

	username := r.FormValue("username")
	password := r.FormValue("password")

	u, err := models.GetUser("username", username)
	if err != nil {
		c.Fail("Unable to find that user")
		return
	}

	ok, err := u.Login(password)
	if err != nil {
		c.Fail("Authentication error")
	}

	if !ok {
		c.Fail("That was not a matching password")
		return
	}

	c.OKWithData("Successfully logged in that user", u)
}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	c := communicator.New(w)

	username := r.FormValue("username")
	password := r.FormValue("password")
	email := r.FormValue("email")

	u, err := models.RegisterUser(username, password, email)
	if err != nil {
		fmt.Println(err)
		c.Fail("Could not register that user")
		return
	}

	c.OKWithData("Successfully registered that user", u)
}
