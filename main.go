package main

import (
	"fmt"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	registerHandlers()
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))

}

func registerHandlers() {
	router := mux.NewRouter()

	router.Methods("GET").Path("/login/").Handler(appHandler(testHandler))

	http.Handle("/", handlers.CombinedLoggingHandler(os.Stderr, router))

}

func testHandler(w http.ResponseWriter, r *http.Request) *appError {
	_, err := fmt.Fprint(w, "TEST PAGE")
	if err != nil {
		return appErrorf(err, "could not write: %v", err)
	}
	return nil
}

type appHandler func(http.ResponseWriter, *http.Request) *appError

type appError struct {
	Error   error
	Message string
	Code    int
}

func (fn appHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if e := fn(w, r); e != nil {
		log.Printf("Handler error: status code %d, message %s, underlying err %#v", e.Code, e.Message, e.Error)
		http.Error(w, e.Message, e.Code)
	}
}

func appErrorf(err error, format string, v ...interface{}) *appError {
	return &appError{
		Error:   err,
		Message: fmt.Sprintf(format, v...),
		Code:    500,
	}
}
