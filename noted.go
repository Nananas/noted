package main

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/nananas/noted/security"
)

type TNote struct {
	Notes []Note
	Prev  string
}

type Note struct {
	Name    string
	ModTime string
}

type Notebook struct {
	Text string
	Name string
}

var SALT string
var VERSION string

func main() {
	log.SetFlags(log.Lshortfile)

	f, err := os.OpenFile("logfile.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Println(err)
	}
	defer f.Close()

	log.SetOutput(f)

	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "-d", "--debug":
			fmt.Println("Starting in debug mode")
			log.SetFlags(log.Lshortfile)
		case "-V", "--version":
			fmt.Println("Noted version " + VERSION)
			return
		case "-h", "--help":
			printHelp()
			return
		}
	}

	security.LoadConfig()
	security.SetSalt(SALT)

	http.HandleFunc("/favicon.ico", handlerStatic)

	// No session required
	//
	http.HandleFunc("/login", handlerLogin)
	http.HandleFunc("/POST_login", handlerPostLogin)

	// Session required
	//
	http.HandleFunc("/", MSession(handler))
	http.HandleFunc("/notes", MSession(handlerNotes))
	http.HandleFunc("/N/", MSession(handlerEditNote))
	http.HandleFunc("/POST_save", MSession(handlerPostSave))

	err = http.ListenAndServe(":4444", nil)
	if err != nil {
		log.Fatal(err)
	}
}

// Middleware: continue if valid session
// 				else redirect to login
//
func MSession(f func(rw http.ResponseWriter, req *http.Request)) http.HandlerFunc {
	return func(rw http.ResponseWriter, req *http.Request) {
		if !security.ValidSession(req) {
			http.Redirect(rw, req, "/login", http.StatusSeeOther)
			return
		}

		// update session if valid
		//
		security.UpdateSession(rw, req)

		f(rw, req)
	}
}

func handlerStatic(rw http.ResponseWriter, req *http.Request) {
	http.NotFound(rw, req)
}

func handler(rw http.ResponseWriter, req *http.Request) {
	http.Redirect(rw, req, "/N/", http.StatusSeeOther)
}

func handlerLogin(rw http.ResponseWriter, req *http.Request) {
	// serve login html file
	http.ServeFile(rw, req, "./html/login.html")
}

func handlerNotes(rw http.ResponseWriter, req *http.Request) {

	infos, err := ioutil.ReadDir("./notes")
	if err != nil {
		// dir does not exist
		err := os.Mkdir("./notes", 0664)
		if err != nil {
			log.Println(err)
			return
		}
	}

	notes := TNote{}

	if prev, ok := req.URL.Query()["prev"]; ok {
		notes.Prev = prev[0]
	}

	for _, i := range infos {
		notes.Notes = append(notes.Notes, Note{
			Name:    i.Name(),
			ModTime: i.ModTime().Format("02-01-2006"),
		})
		// notes = notes + i.Name() + "</br>"
	}

	template, err := template.New("notes.html").ParseFiles("./html/notes.html")
	if err != nil {
		log.Println(err)
		return
	}

	// buffer := bytes.NewBuffer(nil)
	err = template.Execute(rw, notes)
	if err != nil {
		log.Println(err)
	}
}

func handlerEditNote(rw http.ResponseWriter, req *http.Request) {

	notebook := Notebook{}

	if len(req.URL.Path[3:]) != 0 {

		// check if file exists
		b, err := ioutil.ReadFile("./notes/" + req.URL.Path[3:])
		if err != nil {
			log.Println(err)
		} else {
			notebook.Text = string(b)
			notebook.Name = req.URL.Path[3:]
		}
	}

	template, err := template.New("notebook.html").ParseFiles("./html/notebook.html")
	if err != nil {
		log.Println(err)
		return
	}

	// buffer := bytes.NewBuffer(nil)
	err = template.Execute(rw, notebook)
	if err != nil {
		log.Println(err)
	}
	// serve notebook files
	// http.ServeFile(rw, req, "./html/notebook.html")
}

func handlerPostLogin(rw http.ResponseWriter, req *http.Request) {

	if req.Method != "POST" {
		http.Error(rw, "Invalid method.", http.StatusMethodNotAllowed)
	}

	user := req.PostFormValue("user")
	pass := req.PostFormValue("pass")

	if security.Validate(user, pass) {

		log.Println("starting valid session for ", user)

		security.StartSession(rw, user)

		http.Redirect(rw, req, "/", http.StatusSeeOther)

	} else {
		fmt.Fprint(rw, "Username or password incorrect.")
		log.Println("Login attempt: ", user, pass, req.RemoteAddr)
	}

}

func handlerPostSave(rw http.ResponseWriter, req *http.Request) {

	if req.Method != "POST" {
		http.Error(rw, "Invalid method.", http.StatusMethodNotAllowed)
		return
	}

	name := req.PostFormValue("name")
	text := req.PostFormValue("text")

	// remove bad characters
	name = strings.Replace(name, "/", "", -1)
	name = strings.Replace(name, "\\", "", -1)

	// save note as file
	err := ioutil.WriteFile(filepath.Join("./notes", name), []byte(text), 0644)
	if err != nil {
		log.Println(err)
		rw.Write([]byte(err.Error()))
		return
	}

	// redirect to new note
	http.Redirect(rw, req, "/", http.StatusSeeOther)

}

func isAny(b rune, any ...rune) bool {
	for _, a := range any {
		if b == a {
			return true
		}
	}

	return false
}

func printHelp() {
	help := `
Usage:
	./noted [options]

Options:
	-h --help 		print this help message
	-V --version	print version info
	-d --debug		print debug info to stdout instead of logfile.log
`
	fmt.Print(help)
}
