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
	router.HandleFunc("/classes/{class_id}/workshops", restrict.R(CreateWorkshopHandler)).Methods("POST")
	router.HandleFunc("/classes/{class_id}/workshops/{workshop_id}/pages", restrict.R(CreateWorkshopPageHandler)).Methods("POST")
	router.HandleFunc("/classes/{class_id}/workshops/{workshop_id}/pages", restrict.R(GetWorkshopPagesHandler)).Methods("GET")
	router.HandleFunc("/classes/{class_id}/workshops/{workshop_id}/pages/{page_id}/edit", restrict.R(EditPageHandler)).Methods("POST")
	router.HandleFunc("/classes/{class_id}/workshops/{workshop_id}", restrict.R(GetWorkshopHandler)).Methods("GET")
	router.HandleFunc("/classes/{class_id}/workshops/{workshop_id}", restrict.R(DeleteWorkshopHandler)).Methods("DELETE")
	router.HandleFunc("/classes/{class_id}/workshops", restrict.R(GetWorkshopsHandler)).Methods("GET")
	router.HandleFunc("/classes/{class_id}/students", restrict.R(GetStudentHandler)).Methods("GET")
	router.HandleFunc("/classes/{class_id}/image", restrict.R(SetClassImageHandler)).Methods("POST")
	router.HandleFunc("/classes/{class_id}/admin/students", restrict.R(AddUserToClassHandler)).Methods("POST")

	n := negroni.New()
	n.Use(negroni.NewStatic(http.Dir("static")))
	n.Use(negroni.NewLogger())
	n.UseHandler(router)

	models.InitDB(*dbFile)

	http.ListenAndServe("localhost:8080", n)
}

func DeleteWorkshopHandler(w http.ResponseWriter, r *http.Request, t *jwt.Token) {
	c := communicator.New(w)
	vars := mux.Vars(r)

	u, err := getUserFromToken(t)
	if err != nil {
		c.Fail("COuld not get user")
		return
	}

	s, class, err := getClassFromString(u, vars["class_id"])
	if err != nil {
		c.Fail("Could not get class info")
		return
	}

	if !s.IsAdmin() {
		c.Error("You are not an admin. You cannot delete this!")
		return
	}

	wid, err := strconv.ParseInt(vars["workshop_id"], 10, 64)
	if err != nil {
		c.Fail("Not a valid workshop id")
		return
	}

	work, err := class.Workshop(wid)
	if err != nil {
		c.Fail("Could not get workshop")
		return
	}

	err = work.Delete()
	if err != nil {
		c.Fail("Could not delete workshop")
		return
	}

	c.OK("Deleted successfully")
}

func EditPageHandler(w http.ResponseWriter, r *http.Request, t *jwt.Token) {
	c := communicator.New(w)
	vars := mux.Vars(r)

	u, err := getUserFromToken(t)
	if err != nil {
		c.Fail("COuld not get user")
		return
	}

	s, class, err := getClassFromString(u, vars["class_id"])
	if err != nil {
		c.Fail("Could not get class info")
		return
	}

	if !s.IsAdmin() {
		c.Error("You are not an admin. You cannot edit this!")
		return
	}

	wid, err := strconv.ParseInt(vars["workshop_id"], 10, 64)
	if err != nil {
		c.Fail("Not a valid workshop id")
		return
	}

	work, err := class.Workshop(wid)
	if err != nil {
		c.Fail("Could not get workshop")
		return
	}

	pid, err := strconv.ParseInt(vars["page_id"], 10, 64)
	if err != nil {
		c.Fail("Coudl not get page id")
		return
	}

	p, err := work.Page(pid)
	if err != nil {
		c.Fail("COuld not get page")
		return
	}

	err = p.Edit(r.FormValue("title"), r.FormValue("contents"))
	if err != nil {
		c.Fail("Could not edit page")
		return
	}

	c.OKWithData("here is your page", p)
}

func GetWorkshopHandler(w http.ResponseWriter, r *http.Request, t *jwt.Token) {
	c := communicator.New(w)
	vars := mux.Vars(r)

	u, err := getUserFromToken(t)
	if err != nil {
		c.Fail("COuld not get user")
		return
	}

	_, class, err := getClassFromString(u, vars["class_id"])
	if err != nil {
		c.Fail("Could not get class info")
		return
	}

	wid, err := strconv.ParseInt(vars["workshop_id"], 10, 64)
	if err != nil {
		c.Fail("Not a valid workshop id")
		return
	}

	work, err := class.Workshop(wid)
	if err != nil {
		c.Fail("Could not get workshop")
		return
	}

	c.OKWithData("Hre is your workshop", work)
}

func GetWorkshopPagesHandler(w http.ResponseWriter, r *http.Request, t *jwt.Token) {
	c := communicator.New(w)
	vars := mux.Vars(r)

	u, err := getUserFromToken(t)
	if err != nil {
		c.Fail("COuld not get user")
		return
	}

	_, class, err := getClassFromString(u, vars["class_id"])
	if err != nil {
		c.Fail("Could not get class info")
		return
	}

	wid, err := strconv.ParseInt(vars["workshop_id"], 10, 64)
	if err != nil {
		c.Fail("Not a valid workshop id")
		return
	}

	work, err := class.Workshop(wid)
	if err != nil {
		c.Fail("Could not get workshop")
		return
	}

	ps, err := work.Pages()
	if err != nil {
		c.Fail("Could not get pages")
		return
	}

	c.OKWithData("Here are your pages", ps)
}

func CreateWorkshopPageHandler(w http.ResponseWriter, r *http.Request, t *jwt.Token) {
	c := communicator.New(w)
	vars := mux.Vars(r)

	u, err := getUserFromToken(t)
	if err != nil {
		c.Fail("COuld not get user")
		return
	}

	_, class, err := getClassFromString(u, vars["class_id"])
	if err != nil {
		c.Fail("Could not get class info")
		return
	}

	wid, err := strconv.ParseInt(vars["workshop_id"], 10, 64)
	if err != nil {
		c.Fail("Not a valid workshop id")
		return
	}

	work, err := class.Workshop(wid)
	if err != nil {
		c.Fail("Could not get workshop")
		return
	}

	title := r.FormValue("title")
	content := r.FormValue("content")

	p, err := work.CreatePage(title, content)
	if err != nil {
		c.Fail("Could not create page")
		return
	}

	c.OKWithData("Here is your page", p)
}

func SetClassImageHandler(w http.ResponseWriter, r *http.Request, t *jwt.Token) {
	c := communicator.New(w)
	vars := mux.Vars(r)

	u, err := getUserFromToken(t)
	if err != nil {
		c.Fail("Could not get user")
		return
	}

	s, class, err := getClassFromString(u, vars["class_id"])
	if err != nil {
		c.Fail("Could not get class")
		return
	}

	if !s.IsAdmin() {
		c.Error("You are not an admin!")
		return
	}

	err = class.SetImage(r.FormValue("image_url"))
	if err != nil {
		c.Fail("Could not set image url")
		return
	}

	c.OK("Done!")
}

func AddUserToClassHandler(w http.ResponseWriter, r *http.Request, t *jwt.Token) {
	c := communicator.New(w)
	vars := mux.Vars(r)

	u, err := getUserFromToken(t)
	if err != nil {
		c.Fail("COuld not get user from supplied token")
		return
	}

	username := r.FormValue("user")

	inv, err := models.GetUser("username", username)
	if err != nil {
		c.Fail("Could not get user")
		return
	}

	s, class, err := getClassFromString(u, vars["class_id"])
	if err != nil {
		c.Fail("Could not get class")
		return
	}

	if !s.IsAdmin() {
		c.Fail("User is not admin")
		return
	}

	_, err = class.Invite(inv)
	if err != nil {
		c.Fail("Error inviting user")
		return
	}

	c.OK("Invited user")
}

func GetStudentHandler(w http.ResponseWriter, r *http.Request, t *jwt.Token) {
	c := communicator.New(w)
	vars := mux.Vars(r)

	u, err := getUserFromToken(t)
	if err != nil {
		c.Fail("Could not get user in token")
		return
	}

	s, _, err := getClassFromString(u, vars["class_id"])
	if err != nil {
		c.Fail("Unable to get that class")
		return
	}

	c.OKWithData("Here is your data", s)
}

func GetWorkshopsHandler(w http.ResponseWriter, r *http.Request, t *jwt.Token) {
	c := communicator.New(w)

	vars := mux.Vars(r)
	u, err := getUserFromToken(t)
	if err != nil {
		c.Fail("Unable to get user from token")
		return
	}

	s, _, err := getClassFromString(u, vars["class_id"])
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

	as, err := s.Workshops(tm)
	if err != nil {
		c.Fail("Could not get workshop")
		return
	}

	c.OKWithData("Here are your workshops", as)
}

func CreateWorkshopHandler(w http.ResponseWriter, r *http.Request, t *jwt.Token) {
	c := communicator.New(w)

	name := r.FormValue("name")
	description := r.FormValue("description")
	explanation := r.FormValue("explanation")

	if name == "" || description == "" || explanation == "" {
		c.Fail("Invalid data")
		return
	}

	vars := mux.Vars(r)

	u, err := getUserFromToken(t)
	if err != nil {
		c.Fail("Unable to get user from token")
		return
	}

	s, _, err := getClassFromString(u, vars["class_id"])
	if err != nil {
		c.Fail("Could not get class info")
		return
	}

	a, err := s.CreateWorkshop(name, description, explanation)
	if err != nil {
		c.Fail("Could not create workshop")
		return
	}

	c.OKWithData("Here is your workshops", a)
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

	err = class.SetImage("imgs/home.png")
	if err != nil {
		c.Fail("Could not set home class image")
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

func getClassFromString(u models.User, stringID string) (models.Student, models.Class, error) {
	var (
		s models.Student
		c models.Class
	)

	idI, err := strconv.Atoi(stringID) // ugly variable names ahead.
	if err != nil {
		return s, c, err
	}

	id := int64(idI)

	s, c, err = u.Class(id)

	return s, c, err
}
