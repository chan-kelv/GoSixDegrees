package routes

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func Init() {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", Index)
	log.Fatal(http.ListenAndServe(":8080", router))
}

//Says hello! used as test route
func Index(w http.ResponseWriter, req *http.Request) {
	fmt.Println("Hello there")
	fmt.Fprintf(w, "Hello there!")
}
