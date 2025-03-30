package server

import (
	"crypto/rsa"
	"crypto/tls"
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"slices"
	"strconv"
	"time"

	token "github.com/Festivals-App/festivals-identity-server/auth"
	"github.com/Festivals-App/festivals-identity-server/server/config"
	"github.com/Festivals-App/festivals-identity-server/server/database"
	"github.com/Festivals-App/festivals-identity-server/server/handler"
	festivalspki "github.com/Festivals-App/festivals-pki"
	servertools "github.com/Festivals-App/festivals-server-tools"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-jwt/jwt/v5"
	"github.com/rs/zerolog/log"
)

// Server has router and db instances
type Server struct {
	Router    *chi.Mux
	DB        *sql.DB
	Config    *config.Config
	TLSConfig *tls.Config
	Auth      *token.AuthService
	Validator *token.ValidationService
}

func NewServer(config *config.Config) *Server {
	server := &Server{}
	server.initialize(config)
	return server
}

// Initialize the server with predefined configuration
func (s *Server) initialize(config *config.Config) {

	s.Config = config
	s.Router = chi.NewRouter()

	s.setDatabase()
	s.setTLSHandling()
	s.setIdentityService()
	s.setMiddleware()
	s.setRoutes()
}

func (s *Server) setDatabase() {

	dbURI := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True",
		s.Config.DB.Username,
		s.Config.DB.Password,
		s.Config.DB.Host,
		s.Config.DB.Port,
		s.Config.DB.Name,
		s.Config.DB.Charset,
	)

	db, err := sql.Open(s.Config.DB.Dialect, dbURI)

	if err != nil {
		log.Fatal().Err(err).Msg("Failed to open database handle.")
	}

	err = db.Ping()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to connect to database.")
	}

	db.SetConnMaxIdleTime(time.Minute * 1)
	db.SetConnMaxLifetime(time.Minute * 5)

	s.DB = db
}

func (s *Server) setTLSHandling() {

	tlsConfig, err := festivalspki.NewServerTLSConfig(s.Config.TLSCert, s.Config.TLSKey, s.Config.TLSRootCert)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to set TLS handling")
	}
	s.TLSConfig = tlsConfig
}

func (s *Server) setIdentityService() {

	s.Auth = token.NewAuthService(s.Config.AccessTokenPrivateKeyPath, s.Config.AccessTokenPublicKeyPath, s.Config.JwtExpiration, s.Config.ServiceBindHost)
	s.Validator = newLocalValidationService(s.Config.AccessTokenPublicKeyPath)
}

func (s *Server) setMiddleware() {

	// tell the router which middleware to use
	s.Router.Use(
		// used to log the request to the console
		servertools.Middleware(servertools.TraceLogger(s.Config.TraceLog)),
		// tries to recover after panics
		middleware.Recoverer,
	)
}

// setRouters sets the all required routers
func (s *Server) setRoutes() {

	s.Router.Get("/version", s.handleRequest(handler.GetVersion))
	s.Router.Get("/info", s.handleRequest(handler.GetInfo))
	s.Router.Get("/health", s.handleRequest(handler.GetHealth))

	s.Router.Post("/update", s.handleRequest(handler.MakeUpdate))
	s.Router.Get("/log", s.handleRequest(handler.GetLog))
	s.Router.Get("/log/trace", s.handleRequest(handler.GetTraceLog))

	s.Router.Post("/users/signup", s.handleAPIRequest(handler.Signup))
	s.Router.Get("/users/login", s.handleAPIRequest(handler.Login))
	s.Router.Get("/users/refresh", s.handleRequest(handler.Refresh))
	s.Router.Get("/users", s.handleRequest(handler.GetUsers))
	s.Router.Post("/users/{objectID}/change-password", s.handleRequest(handler.ChangePassword))
	s.Router.Post("/users/{objectID}/suspend", s.handleRequest(handler.SuspendUser))
	s.Router.Post("/users/{objectID}/role/{resourceID}", s.handleRequest(handler.SetUserRole))

	s.Router.Post("/users/{objectID}/festival/{resourceID}", s.handleServiceRequest(handler.SetFestivalForUser))
	s.Router.Post("/users/{objectID}/artist/{resourceID}", s.handleServiceRequest(handler.SetArtistForUser))
	s.Router.Post("/users/{objectID}/location/{resourceID}", s.handleServiceRequest(handler.SetLocationForUser))
	s.Router.Post("/users/{objectID}/event/{resourceID}", s.handleServiceRequest(handler.SetEventForUser))
	s.Router.Post("/users/{objectID}/link/{resourceID}", s.handleServiceRequest(handler.SetLinkForUser))
	s.Router.Post("/users/{objectID}/image/{resourceID}", s.handleServiceRequest(handler.SetImageForUser))
	s.Router.Post("/users/{objectID}/place/{resourceID}", s.handleServiceRequest(handler.SetPlaceForUser))
	s.Router.Post("/users/{objectID}/tag/{resourceID}", s.handleServiceRequest(handler.SetTagForUser))

	s.Router.Delete("/users/{objectID}/festival/{resourceID}", s.handleServiceRequest(handler.RemoveFestivalForUser))
	s.Router.Delete("/users/{objectID}/artist/{resourceID}", s.handleServiceRequest(handler.RemoveArtistForUser))
	s.Router.Delete("/users/{objectID}/location/{resourceID}", s.handleServiceRequest(handler.RemoveLocationForUser))
	s.Router.Delete("/users/{objectID}/event/{resourceID}", s.handleServiceRequest(handler.RemoveEventForUser))
	s.Router.Delete("/users/{objectID}/link/{resourceID}", s.handleServiceRequest(handler.RemoveLinkForUser))
	s.Router.Delete("/users/{objectID}/image/{resourceID}", s.handleServiceRequest(handler.RemoveImageForUser))
	s.Router.Delete("/users/{objectID}/place/{resourceID}", s.handleServiceRequest(handler.RemovePlaceForUser))
	s.Router.Delete("/users/{objectID}/tag/{resourceID}", s.handleServiceRequest(handler.RemoveTagForUser))

	s.Router.Get("/validation-key", s.handleServiceRequest(handler.GetValidationKey))

	s.Router.Get("/api-keys", s.handleServiceRequest(handler.GetAPIKeys))
	s.Router.Post("/api-keys", s.handleRequest(handler.AddAPIKey))
	s.Router.Delete("/api-keys", s.handleRequest(handler.DeleteAPIKey))

	s.Router.Get("/service-keys", s.handleServiceRequest(handler.GetServiceKeys))
	s.Router.Post("/service-keys", s.handleRequest(handler.AddServiceKey))
	s.Router.Delete("/service-keys", s.handleRequest(handler.DeleteServiceKey))
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

type JWTAuthenticatedHandlerFunction func(auth *token.AuthService, claims *token.UserClaims, db *sql.DB, w http.ResponseWriter, r *http.Request)

func (s *Server) handleRequest(requestHandler JWTAuthenticatedHandlerFunction) http.HandlerFunc {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		claims := token.GetValidClaims(r, s.Validator)
		if claims == nil {
			servertools.UnauthorizedResponse(w)
			return
		}
		requestHandler(s.Auth, claims, s.DB, w, r)
	})
}

type APIKeyAuthenticatedHandlerFunction func(auth *token.AuthService, db *sql.DB, w http.ResponseWriter, r *http.Request)

func (s *Server) handleAPIRequest(requestHandler APIKeyAuthenticatedHandlerFunction) http.HandlerFunc {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		apikey := token.GetAPIToken(r)
		allAPIKeys, err := database.GetAllAPIKeys(s.DB)
		if err != nil {
			log.Error().Msg("failed to load API keys from database")
			servertools.RespondError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
			return
		}
		if !slices.Contains(getAPIKeyValues(allAPIKeys), apikey) {
			servertools.UnauthorizedResponse(w)
			return
		}
		requestHandler(s.Auth, s.DB, w, r)
	})
}

type ServiceKeyAuthenticatedHandlerFunction func(auth *token.AuthService, db *sql.DB, w http.ResponseWriter, r *http.Request)

func (s *Server) handleServiceRequest(requestHandler ServiceKeyAuthenticatedHandlerFunction) http.HandlerFunc {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		servicekey := token.GetServiceToken(r)
		if servicekey == "" {
			claims := token.GetValidClaims(r, s.Validator)
			if claims != nil && claims.UserRole == token.ADMIN {
				requestHandler(s.Auth, s.DB, w, r)
				return
			}
			servertools.UnauthorizedResponse(w)
			return
		}
		allServiceKeys, err := database.GetAllServiceKeys(s.DB)
		if err != nil {
			log.Error().Msg("failed to load servive keys from database")
			servertools.RespondError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
			return
		}
		if !slices.Contains(getServiceKeyValues(allServiceKeys), servicekey) {
			servertools.UnauthorizedResponse(w)
			return
		}
		requestHandler(s.Auth, s.DB, w, r)
	})
}

func getServiceKeyValues(keys []token.ServiceKey) []string {
	var data []string
	for _, key := range keys {
		data = append(data, key.Key)
	}
	return data
}

func getAPIKeyValues(keys []token.APIKey) []string {
	var data []string
	for _, key := range keys {
		data = append(data, key.Key)
	}
	return data
}

func newLocalValidationService(publickey string) *token.ValidationService {

	var verifyKey *rsa.PublicKey = nil
	verifyBytes, err := os.ReadFile(publickey)
	if err != nil {
		log.Fatal().Err(err).Msg("Unable to read public auth key.")
	}
	verifyKey, err = jwt.ParseRSAPublicKeyFromPEM(verifyBytes)
	if err != nil {
		log.Fatal().Err(err).Msg("Unable to parse public auth key.")
	}

	return &token.ValidationService{Key: verifyKey, APIKeys: nil, ServiceKeys: nil, Client: nil, Endpoint: ""}
}
