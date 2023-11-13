package server

import (
	"crypto/tls"
	"database/sql"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/Festivals-App/festivals-gateway/server/logger"
	"github.com/Festivals-App/festivals-identity-server/server/config"
	"github.com/Festivals-App/festivals-identity-server/server/handler"
	festivalspki "github.com/Festivals-App/festivals-pki"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	_ "github.com/go-sql-driver/mysql"
	"github.com/rs/zerolog/log"
)

// Server has router and db instances
type Server struct {
	Router    *chi.Mux
	DB        *sql.DB
	Config    *config.Config
	TLSConfig *tls.Config
}

func NewServer(config *config.Config) *Server {
	server := &Server{}
	server.Initialize(config)
	return server
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
	db, _ := sql.Open(config.DB.Dialect, dbURI)

	//if err != nil {
	//log.Fatal().Err((err)).Msg("Server initialize: could not connect to database")
	//}

	s.DB = db
	s.Router = chi.NewRouter()
	s.Config = config

	s.setTLSHandling()
	s.setMiddleware()
	s.setRoutes()
}

func (s *Server) setTLSHandling() {

	tlsConfig := &tls.Config{
		ClientAuth:     tls.RequireAndVerifyClientCert,
		GetCertificate: festivalspki.LoadServerCertificateHandler(s.Config.TLSCert, s.Config.TLSKey, s.Config.TLSRootCert),
	}
	s.TLSConfig = tlsConfig
}

func (s *Server) setMiddleware() {

	// tell the router which middleware to use
	s.Router.Use(
		// used to log the request to the console
		logger.Middleware(logger.TraceLogger("/var/log/festivals-identity-server/trace.log")),
		// tries to recover after panics (?)
		middleware.Recoverer,
	)
}

// setRouters sets the all required routers
func (s *Server) setRoutes() {

	s.Router.Get("/version", s.handleRequestWithoutValidation(handler.GetVersion))
	s.Router.Get("/info", s.handleRequestWithoutValidation(handler.GetInfo))
	s.Router.Get("/health", s.handleRequestWithoutValidation(handler.GetHealth))

	s.Router.Post("/update", s.handleRequest(handler.MakeUpdate))
	s.Router.Get("/log", s.handleRequest(handler.GetLog))
	s.Router.Get("/log/trace", s.handleRequest(handler.GetTraceLog))

	s.Router.Post("/signup", s.handleRequestWithoutValidation(handler.Signup))
	s.Router.Post("/login", s.handleRequestWithoutValidation(handler.Login))

	s.Router.Post("/refresh", s.handleRequest(handler.Signup))
	s.Router.Post("/change-password", s.handleRequest(handler.Signup))

	//s.Router.Get("/user", s.handleRequest(handler.GetLog))

	//s.Router.Get("/festivals", s.handleRequest(handler.GetFestivals))
	//s.Router.Get("/festivals/{objectID}", s.handleRequest(handler.GetFestival))
}

func (s *Server) Run(conf *config.Config) {

	server := http.Server{
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      15 * time.Second,
		IdleTimeout:       60 * time.Second,
		ReadHeaderTimeout: 5 * time.Second,
		Addr:              conf.ServiceBindHost + ":" + strconv.Itoa(conf.ServicePort),
		Handler:           s.Router,
		TLSConfig:         s.TLSConfig,
	}

	if err := server.ListenAndServeTLS("", ""); err != nil {
		log.Fatal().Err(err).Str("type", "server").Msg("Failed to run server")
	}
}

// function prototype to inject DB instance in handleRequest()
type RequestHandlerFunction func(db *sql.DB, w http.ResponseWriter, r *http.Request)

func (s *Server) handleRequest(requestHandler RequestHandlerFunction) http.HandlerFunc {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestHandler(s.DB, w, r)
	})
}

func (s *Server) handleRequestWithoutValidation(requestHandler RequestHandlerFunction) http.HandlerFunc {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestHandler(s.DB, w, r)
	})
}
