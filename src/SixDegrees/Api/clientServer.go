package routes

import (
	pages "SixDegrees/Pages"
	node "SixDegrees/Peers"
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func Init() {
	router := mux.NewRouter().StrictSlash(true)
	handleEndpoints(router)

	//Static load the folder (to load the css page)
	router.PathPrefix("/Pages/").Handler(http.StripPrefix("/Pages/", http.FileServer(http.Dir("./Pages/"))))
	log.Fatal(http.ListenAndServe(":8080", router))
}

func handleEndpoints(router *mux.Router) {
	router.HandleFunc("/", launch)
	router.HandleFunc("/crawl", crawl)
}

func crawl(w http.ResponseWriter, req *http.Request) {
	startCrawlTerm := req.FormValue("crawlTerm")
	fmt.Println("Start crawl on:", startCrawlTerm)
	node.MakeTestRpc(startCrawlTerm)
}

func launch(w http.ResponseWriter, req *http.Request) {
	fmt.Println("method:", req.Method) //get request method
	t, err := template.ParseFiles(pages.LaunchPage)
	if err != nil {
		fmt.Println("launch page error:", err)
	}
	t.Execute(w, nil)
}
