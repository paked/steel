package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/codegangsta/negroni"
	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"github.com/paked/gerrycode/communicator"
	"github.com/paked/steel/models"
)

var (
	pkey []byte
)

func main() {
	var err error
	pkey, err = readPrivateKey()
	if err != nil {
		panic(err)
	}

	router := mux.NewRouter()
	router.HandleFunc("/user", GetUserHandler).Methods("GET")
	router.HandleFunc("/user/login", LoginHandler).Methods("POST")
	router.HandleFunc("/user/register", RegisterHandler).Methods("POST")
	router.HandleFunc("/assignments", restrict(CreateAssignmentHandler)).Methods("POST")

	router.Handle("/cake", restrict(GiveCakeHandler))

	n := negroni.New()
	n.Use(negroni.NewStatic(http.Dir("static")))
	n.Use(negroni.NewLogger())
	n.UseHandler(router)

	models.InitDB()

	http.ListenAndServe("localhost:8080", n)
}

func GiveCakeHandler(w http.ResponseWriter, r *http.Request, t *jwt.Token) {
	fmt.Fprintln(w, "*cake*")
}

func GetUserHandler(w http.ResponseWriter, r *http.Request) {
	c := communicator.New(w)
	username := r.FormValue("username")
	sid := r.FormValue("id")

	if sid != "" {
		var id int64
		fmt.Sscanf(sid, "%d", &id)
		fmt.Println(id)

		u, err := models.GetUserByID(id)
		if err != nil {
			c.Fail("Error getting user")
			return
		}

		c.OKWithData("Here is your user: ", u)
		return
	}

	if username != "" {
		u, err := models.GetUser("username", username)
		if err != nil {
			c.Fail("Could not get that username")
			return
		}

		c.OKWithData("Here is your user: ", u)
		return
	}

	restrict(func(w http.ResponseWriter, r *http.Request, t *jwt.Token) {
		c := communicator.New(w)
		id, ok := t.Claims["id"].(int64)
		if !ok {
			c.Fail("Could not get that ID")
			return
		}

		u, err := models.GetUserByID(id)
		if err != nil {
			c.Fail("Error getting user")
			return
		}

		c.OKWithData("Here is your user", u)
	})

	c.Fail("You didn't provide any data :/")
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

	token := jwt.New(jwt.SigningMethodHS256)
	token.Claims["id"] = u.ID
	token.Claims["exp"] = time.Now().Add(time.Hour * 72).Unix()

	ts, err := token.SignedString(pkey)
	if err != nil {
		c.Fail("Failure signing that token!")
		fmt.Println(err)
		return
	}

	c.OKWithData("Successfully logged in that user", ts)
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

func CreateAssignmentHandler(w http.ResponseWriter, r *http.Request, t *jwt.Token) {
	c := communicator.New(w)

	id, ok := t.Claims["id"].(int64)
	if !ok {
		c.Fail("Not a valid ID in token")
		return
	}

	u, err := models.GetUserByID(id)
	if err != nil {
		c.Fail("Unable to get user")
		return
	}

	name := r.FormValue("name")
	description := r.FormValue("description")
	explanation := r.FormValue("explanation")

	if name == "" || description == "" || explanation == "" {
		c.Fail("Not valid name/description/explanation")
		return
	}

	a, err := u.CreateAssignment(name, description, explanation)
	if err != nil {
		c.Fail("Could not create assignment " + err.Error())
		return
	}

	c.OKWithData("Here is the assignment", a)
}

func restrict(fn func(http.ResponseWriter, *http.Request, *jwt.Token)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ts := r.FormValue("access_token")
		c := communicator.New(w)

		token, err := jwt.Parse(ts, func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Unexpected signing method: %v", t.Header["alg"])
			}
			return pkey, nil
		})

		if err != nil {
			c.Fail("You are not using a valid token:" + err.Error())
			fmt.Println(err)
			return
		}

		if !token.Valid {
			c.Fail("Something obscurely weird happened to your token!")
			return
		}

		fn(w, r, token)
	}
}

func readPrivateKey() ([]byte, error) {
	privateKey, e := ioutil.ReadFile("keys/app.rsa")

	return privateKey, e
}
