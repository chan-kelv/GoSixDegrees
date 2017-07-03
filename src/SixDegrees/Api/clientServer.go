package routes

import (
	pages "SixDegrees/Pages"
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func Init() {
	router := mux.NewRouter().StrictSlash(true)
	handleEndpoints(router)
	log.Fatal(http.ListenAndServe(":8080", router))
}

//Says hello! used as test route
func Index(w http.ResponseWriter, req *http.Request) {
	fmt.Println("Hello there")
	fmt.Fprintf(w, "Hello there!")
}

func Login(w http.ResponseWriter, req *http.Request) {
	user := req.FormValue("username")
	pw := req.FormValue("password")

	fmt.Println(user, pw)
}

func launch(w http.ResponseWriter, req *http.Request) {
	fmt.Println("method:", req.Method) //get request method
	// if req.Method == "GET" {
	// t, _ := template.ParseFiles("Pages/launch.html")
	t, _ := template.ParseFiles(pages.LaunchPage)
	t.Execute(w, nil)
	// } else {
	// 	req.ParseForm()
	// 	// logic part of log in
	// 	fmt.Println("username:", req.Form["username"])
	// 	fmt.Println("password:", req.Form["password"])
	// }
}

func handleEndpoints(router *mux.Router) {
	router.HandleFunc("/", Index)
	router.HandleFunc("/launch", launch)
	router.HandleFunc("/login", Login)
}
