package main

import (
	"context"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
	"os/signal"
	"quickly-api/handlers"
	"quickly-api/services"
	"syscall"
	"time"
)

const (
	BindAddr     = ":8080"
	ReadTimeout  = 5
	WriteTimeout = 10
	IdleTimeout  = 120
)

func main() {
	l := log.New(os.Stdout, "Quickly API: ", log.LstdFlags)

	//oauthMiddleware := services.NewOauthClientMiddleware()
	logMiddleware := services.NewLogHTTPMiddleware(l)

	mux := mux.NewRouter().StrictSlash(true)
	//mux.Use(oauthMiddleware.Func())
	mux.Use(logMiddleware.Func())
	mux.Use(services.MiddlewareSetContentTypeJson())

	notebooksRouter := mux.PathPrefix("/notebooks").Subrouter()
	notebooksRouter.HandleFunc("/", handlers.ReadNotebooks).Methods("GET")
	notebooksRouter.HandleFunc("/", handlers.CreateNotebook).Methods("POST")
	notebooksRouter.HandleFunc("/", handlers.UpdateNotebooks).Methods("PUT")
	notebooksRouter.HandleFunc("/", handlers.DeleteNotebooks).Methods("DELETE")
	notebooksRouter.HandleFunc("/{id}", handlers.ReadNotebookById).Methods("GET")
	notebooksRouter.HandleFunc("/{id}", handlers.UpdateNotebookById).Methods("PUT")
	notebooksRouter.HandleFunc("/{id}", handlers.DeleteNotebookById).Methods("DELETE")

	notesRouter := notebooksRouter.PathPrefix("/{notebooksId}/notes").Subrouter()
	notesRouter.HandleFunc("/", handlers.ReadNotes).Methods("GET")
	notesRouter.HandleFunc("/", handlers.CreateNote).Methods("POST")
	notesRouter.HandleFunc("/", handlers.UpdateNotes).Methods("PUT")
	notesRouter.HandleFunc("/", handlers.DeleteNotes).Methods("DELETE")
	notesRouter.HandleFunc("/{id}", handlers.ReadNoteById).Methods("GET")
	notesRouter.HandleFunc("/{id}", handlers.UpdateNoteById).Methods("PUT")
	notesRouter.HandleFunc("/{id}", handlers.DeleteNoteById).Methods("DELETE")

	s := http.Server{
		Addr:         BindAddr,                   // configure the bind address
		Handler:      mux,                        // set the default handler
		ErrorLog:     l,                          // set the logger for the server
		ReadTimeout:  ReadTimeout * time.Second,  // max time to read request from the client
		WriteTimeout: WriteTimeout * time.Second, // max time to write response to the client
		IdleTimeout:  IdleTimeout * time.Second,  // max time for connections using TCP Keep-Alive
	}

	go func() {
		l.Printf("Starting server on %s\n", s.Addr)
		err := s.ListenAndServe()
		if err != nil {
			l.Fatalln(err)
		}
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)
	signal.Notify(sigChan, syscall.SIGTERM)

	sig := <-sigChan
	l.Printf("Received signal %s\n", sig)

	l.Println("Starting graceful shutdown")
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	s.Shutdown(ctx)

}
