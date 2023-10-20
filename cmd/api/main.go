package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"plutus/internal/api"
	"plutus/internal/model"

	"github.com/gorilla/mux"
)

const (
	defaultDsn  = ""
	defaultPort = "8080"
)

type cfg struct {
	Port string
	Dsn  string
}

func parseConfig() (c cfg) {
	dsn := os.Getenv("PSQL_DSN")
	if len(dsn) == 0 {
		dsn = defaultDsn
	}

	port := os.Getenv("SERVER_PORT")
	if len(port) == 0 {
		port = defaultPort
	}

	return cfg{
		Port: port,
		Dsn:  dsn,
	}
}

func main() {
	config := parseConfig()

	// Check setup db pool
	err := model.InitDBPool(config.Dsn)
	if err != nil {
		panic(err)
	}

	// Init table
	repo, err := model.NewDBRepo(context.Background())
	if err != nil {
		panic(err)
	}
	err = repo.Init()
	if err != nil {
		panic(err)
	}
	repo.Close()

	router := mux.NewRouter()

	// Set DbRepository middleware
	router.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			rr := r.WithContext(context.WithValue(r.Context(), api.GetRepo,
				func(r context.Context) (model.DbRepository, error) {
					return model.NewDBRepo(r)
				}))
			next.ServeHTTP(w, rr)
		})
	})

	api.SetRoutes(router)

	slog.Info("Server started!")
	err = http.ListenAndServe(fmt.Sprintf(":%s", config.Port), router)
	if err != nil {
		panic(err)
	}
}
