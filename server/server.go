package server

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/Festivals-App/festivals-identity-server/server/config"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

// Server has router and db instances
type Server struct {
	Router *chi.Mux
	DB     *sql.DB
}

// Initialize the server with predefined configuration
func (s *Server) Initialize(config *config.Config) {

	dbURI := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True",
		config.DB.Username,
		config.DB.Password,
		config.DB.Host,
		config.DB.Port,
		config.DB.Name,
		config.DB.Charset)
	db, err := sql.Open(config.DB.Dialect, dbURI)
	if err != nil {
		log.Fatal("server initialize: could not connect to database")
	}

	s.DB = db
	s.Router = chi.NewRouter()

	s.setMiddleware()
	s.setWalker()
	s.setRoutes(config)
}

func (s *Server) setMiddleware() {
	// tell the router which middleware to use
	s.Router.Use(
		// used to log the request to the console | development
		middleware.Logger,
		// helps to redirect wrong requests (why do one want that?)
		//middleware.RedirectSlashes,
		// tries to recover after panics (?)
		middleware.Recoverer,
	)
}

func (s *Server) setWalker() {

	walkFunc := func(method string, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
		log.Printf("%s %s \n", method, route)
		return nil
	}
	if err := chi.Walk(s.Router, walkFunc); err != nil {
		log.Panicf("Logging err: %s\n", err.Error())
	}
}

// setRouters sets the all required routers
func (s *Server) setRoutes(config *config.Config) {

}

// Run the server on it's router
func (s *Server) Run(host string) {
	log.Fatal(http.ListenAndServe(host, s.Router))
}

// function prototype to inject DB instance in handleRequest()
type RequestHandlerFunction func(db *sql.DB, w http.ResponseWriter, r *http.Request)

// inject DB in handler functions
func (s *Server) handleRequest(handler RequestHandlerFunction) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		handler(s.DB, w, r)
	}
}
