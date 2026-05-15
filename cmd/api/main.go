package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/alexedwards/scs/pgxstore"
	"github.com/alexedwards/scs/v2"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jmarcosnsf/gobank/internal/api"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		panic(err)
	}

	ctx := context.Background()
	pool, err := pgxpool.New(ctx, fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s",
		os.Getenv("GOBANK_DATABASE_USER"),
		os.Getenv("GOBANK_DATABASE_PASSWORD"),
		os.Getenv("GOBANK_DATABASE_HOST"),
		os.Getenv("GOBANK_DATABASE_PORT"),
		os.Getenv("GOBANK_DATABASE_NAME"),
	))
	if err != nil {
		panic(err)
	}

	defer pool.Close()

	s := scs.New()
	s.Store = pgxstore.New(pool)
	s.Lifetime = 24 * time.Hour
	s.Cookie.HttpOnly = true
	s.Cookie.SameSite = http.SameSiteLaxMode

	api := api.Api{
		Router: chi.NewRouter(),
		Sessions: s,
	}

	api.BindRoutes()

	fmt.Println("Starting Server on port :3080")
	if err := http.ListenAndServe("localhost:3080", api.Router); err != nil {
		panic (err)
	}
}
