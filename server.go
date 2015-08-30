package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/codegangsta/negroni"
	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"github.com/paked/configure"
	"github.com/paked/gerrycode/communicator"
	"github.com/paked/restrict"
	"github.com/paked/steel/models"
)

var (
	pkey []byte

	conf   = configure.New()
	dbFile = conf.String("db", "database.db", "path to the db")
)

func init() {
	conf.Use(configure.NewFlag())
	conf.Use(configure.NewJSONFromFile("config.json"))
}

func main() {
	conf.Parse()

	restrict.ReadCryptoKey("keys/app.rsa")

	router := mux.NewRouter()
	router.HandleFunc("/users", GetUserHandler).Methods("GET")
	router.HandleFunc("/users/login", LoginHandler).Methods("POST")
	router.HandleFunc("/users/register", RegisterHandler).Methods("POST")
	router.HandleFunc("/classes", restrict.R(CreateClassHandler)).Methods("POST")
	router.HandleFunc("/classes", restrict.R(GetClassesHandler)).Methods("GET")
	router.HandleFunc("/classes/{class_id}/assignments", restrict.R(CreateAssignmentHandler)).Methods("POST")
	router.HandleFunc("/classes/{class_id}/assignments/due", restrict.R(GetDueAssignments)).Methods("GET")

	n := negroni.New()
	n.Use(negroni.NewStatic(http.Dir("static")))
	n.Use(negroni.NewLogger())
	n.UseHandler(router)

	models.InitDB(*dbFile)

	http.ListenAndServe("localhost:8080", n)
}

func GetDueAssignments(w http.ResponseWriter, r *http.Request, t *jwt.Token) {
	c := communicator.New(w)

	vars := mux.Vars(r)
	cID := vars["class_id"]

	idI, err := strconv.Atoi(cID)
	if err != nil {
		c.Fail("Unable to parse that id")
		return
	}

	u, err := getUserFromToken(t)
	if err != nil {
		c.Fail("Unable to get user from token")
		return
	}

	id := int64(idI)

	s, _, err := u.Class(id)
	if err != nil {
		c.Fail("Could not get class info")
		return
	}
	d, err := time.ParseDuration("168h")
	if err != nil {
		c.Fail("could not parse duration")
		return
	}

	tm := time.Now().Add(d)

	as, err := s.DueAssignments(tm)
	if err != nil {
		c.Fail("Could not get assignments")
		return
	}

	c.OKWithData("Here are your assignments", as)
}

func CreateAssignmentHandler(w http.ResponseWriter, r *http.Request, t *jwt.Token) {
	c := communicator.New(w)

	name := r.FormValue("name")
	description := r.FormValue("description")
	explanation := r.FormValue("explanation")

	if name == "" || description == "" || explanation == "" {
		c.Fail("Invalid data")
		return
	}

	vars := mux.Vars(r)
	cID := vars["class_id"]

	idI, err := strconv.Atoi(cID)
	if err != nil {
		c.Fail("Unable to parse that id")
		return
	}

	id := int64(idI)

	u, err := getUserFromToken(t)
	if err != nil {
		c.Fail("Unable to get user from token")
		return
	}

	s, _, err := u.Class(id)
	if err != nil {
		c.Fail("Could not get class info")
		return
	}

	a, err := s.CreateAssignment(name, description, explanation)
	if err != nil {
		c.Fail("Could not create assignmnet")
		return
	}

	c.OKWithData("Here is your assignment", a)
}

func CreateClassHandler(w http.ResponseWriter, r *http.Request, t *jwt.Token) {
	c := communicator.New(w)

	u, err := getUserFromToken(t)
	if err != nil {
		c.Fail("Unable to get your user")
		return
	}

	name := r.FormValue("name")
	description := r.FormValue("description")

	if name == "" || description == "" {
		c.Fail("Invalide data")
		return
	}

	class, err := u.NewClass(name, description)
	if err != nil {
		c.Fail("Could not create class")
		return
	}

	c.OKWithData("Here is your class", class)
}

func GetClassesHandler(w http.ResponseWriter, r *http.Request, t *jwt.Token) {
	c := communicator.New(w)

	u, err := getUserFromToken(t)
	if err != nil {
		c.Fail("Unable to get your user")
		return
	}

	classes, err := u.Classes()
	if err != nil {
		c.Fail("Unable to get classes")
		return
	}

	c.OKWithData("Here are your classes", classes)
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

	restrict.R(func(w http.ResponseWriter, r *http.Request, t *jwt.Token) {
		fmt.Println("hello?")
		c := communicator.New(w)
		fid, ok := t.Claims["id"].(float64)
		if !ok {
			c.Fail("Could not get that ID")
			return
		}

		id := int64(fid)

		u, err := models.GetUserByID(id)
		if err != nil {
			c.Fail("Error getting user")
			return
		}

		c.OKWithData("Here is your user", u)

	})(w, r)
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

	claims := make(map[string]interface{})
	claims["id"] = u.ID
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix()

	ts, err := restrict.Token(claims)
	if err != nil {
		c.Fail("Failure signing that token!")
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

func getUserFromToken(t *jwt.Token) (models.User, error) {
	fid, ok := t.Claims["id"].(float64)
	if !ok {
		return models.User{}, errors.New("Could not get user from token")
	}

	id := int64(fid)

	return models.GetUserByID(id)
}
