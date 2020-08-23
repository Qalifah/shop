package routes

import (
	"encoding/json"
	"net/http"
	"time"
	"github.com/Qalifah/shop/auth"
	"github.com/Qalifah/shop/graphquery"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"github.com/swaggo/http-swagger"
	_ "github.com/swaggo/http-swagger/example/go-chi/docs" // hey! :)
)

//Routers contains the relationship between the routes and the controllers
func Routers() *chi.Mux {
	router := chi.NewRouter()
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
    	AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token", "Accept-Encoding", "Host"},
    	ExposedHeaders:   []string{"Link"},
    	AllowCredentials: false,
    	MaxAge:           300,
	}))
	router.Use(CommonMiddleware, middleware.RequestID, middleware.Logger, middleware.RealIP, middleware.Recoverer, middleware.Timeout(60 * time.Second))
	router.Route("/auth", func(router chi.Router) {
		router.Post("/register", auth.Register)
		router.Post("/login", auth.Login)
		router.Post("/logout", auth.Logout)
		router.Post("/refresh/token", auth.Refresh)
	})
	router.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("http://localhost:1323/swagger/doc.json"), //The url pointing to API definition
	))
	router.With(auth.TokenAuthMiddleware).Post("/graphql", func(w http.ResponseWriter, r *http.Request) {
		query := []byte{}
		r.Body.Read(query)
		result := graphquery.ExecuteQuery(string(query), graphquery.Schema)
		json.NewEncoder(w).Encode(result)

	})
	return router
}

// CommonMiddleware adds some essential headers in our response
func CommonMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, Access-Control-Request-Headers, Access-Control-Request-Method, Connection, Host, Origin, User-Agent, Referer, Cache-Control, X-header")
		next.ServeHTTP(w, r)
	})
}
	