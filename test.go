package main

import (
	"log"
	"net/http"
	"github.com/gorilla/mux"
)

func loggingMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Do stuff here
		log.Println(r.RequestURI)
		w.Write([]byte("this is midleware func"))
        // Call the next handler, which can be another middleware in the chain, or the final handler.
        next.ServeHTTP(w, r)
    })
}
func middleware(w http.ResponseWriter, r *http.Request) {
	log.Println("with finalHandler")
	w.Write([]byte("with  middleware"))
}

func final(w http.ResponseWriter, r *http.Request) {
	log.Println("without Executing finalHandler  ")
	w.Write([]byte("without Executing finalHandler  "))
}
func ProductsHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("this is subrouter "))
}

func main() {  
	r := mux.NewRouter()
	s := r.PathPrefix("/api").Subrouter()
	s.HandleFunc("/product", ProductsHandler)
	r.HandleFunc("/final", final)
	s.HandleFunc("/", middleware)
	r.Use(loggingMiddleware)
	log.Fatal(http.ListenAndServe(":8000", r))

}

//
/*
long report
short report
update order
add order
changekitchenstatus
raiseissue
DisplayTokenList
UpdateTakeawayBill
*/
