package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/Skb08/project01/internal/database"

	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"

	_ "github.com/lib/pq"
)

type apiConfig struct {
	DB *database.Queries
}

func main() {
	godotenv.Load(".env")

	portString := os.Getenv("PORT")
	if portString == "" {
		log.Fatal("PORT environment variable was not set")
	}

	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL environment variable is not set")
	}

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal(err)
	}
	dbQueries := database.New(db)

	apiCfg := apiConfig{
		DB: dbQueries,
	}

	router := chi.NewRouter()

	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	v1Router := chi.NewRouter()
	
	v1Router.Post("/users", apiCfg.handlerUsersCreate)
	v1Router.Get("/users", apiCfg.middlewareAuth(apiCfg.handlerUsersGet))
	
	v1Router.Post("/feeds", apiCfg.middlewareAuth(apiCfg.handlerFeedCreate))
	v1Router.Get("/feeds", apiCfg.handlerGetFeeds)

	v1Router.Get("/feed_follows", apiCfg.middlewareAuth(apiCfg.handlerFeedFollowsGet))
	v1Router.Post("/feed_follows", apiCfg.middlewareAuth(apiCfg.handlerFeedFollowCreate))
	v1Router.Delete("/feed_follows/{feedFollowID}", apiCfg.middlewareAuth(apiCfg.handlerFeedFollowDelete))

	v1Router.Get("/healthz", handlerReadiness)
	v1Router.Get("/err", handlerErr)

	router.Mount("/v1", v1Router)
	srv := &http.Server{
		Addr:    ":" + portString,
		Handler: router,
	}

	log.Printf("Serving on port: %s\n", portString)
	log.Fatal(srv.ListenAndServe())
}
