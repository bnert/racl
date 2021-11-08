package api

import (
    "context"
    "net/http"
    "github.com/go-chi/chi/v5"
    "github.com/go-chi/chi/v5/middleware"
    "github.com/jackc/pgx/v4/pgxpool"
)

type Mw = func(http.ResponseWriter, *http.Request) http.HandlerFunc

func attachDbPool(pool *pgxpool.Pool) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            ctx := context.WithValue(r.Context(), "db", pool)
            next.ServeHTTP(w, r.WithContext(ctx))
        })
    }
}

func ctxDbConn(r *http.Request) (*pgxpool.Conn, error) {
    pool := r.Context().Value("db").(*pgxpool.Pool)

    return pool.Acquire(context.Background())
}

func Router(pool *pgxpool.Pool) chi.Router {
    r := chi.NewRouter()

    r.Use(middleware.RequestID)
    r.Use(middleware.RealIP)
    r.Use(middleware.Logger)

    r.Use(attachDbPool(pool))

    r.Route("/acl", func(r chi.Router) {
        r.Get("/{entityId}", handler(getAcl))
        r.Post("/", handler(createAcl))
        r.Put("/{entityId}", handler(updateAcl))
        r.Delete("/{entityId}", handler(deleteAcl))
    })

    return r
}

