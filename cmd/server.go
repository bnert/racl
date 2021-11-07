package cmd

import (
    "context"
    "log"
    "embed"
    "fmt"
    "net/http"
    "path"
    "os"
    "github.com/spf13/cobra"

    "github.com/bnert/racl/api"
    "github.com/jackc/pgx/v4/pgxpool"
)

var (
    listenAddr string
)

func ensureGrootUser(conn *pgxpool.Pool) {
    
}

func ensureTablesCreated(conn *pgxpool.Pool, sqlEmbeds *embed.FS) {
    basePath := "sql/schema"

    dir, err := sqlEmbeds.ReadDir(basePath)
    if err != nil {
        fmt.Println("Cannot open file", basePath)
        return
    }

    fmt.Println("Ensuring uuid-ossp")
    conn.Exec(context.Background(), "create extension if not exists \"uuid-ossp\"")

    for _, v := range dir {
        p := path.Join(basePath, v.Name())
        f, err := sqlEmbeds.ReadFile(p)
        if err != nil {
            fmt.Println("Cannot open file:", p)
        }

        fmt.Println("Migrating:", v.Name())
        conn.Exec(context.Background(), string(f))
    }
}

func registerServerCmd(sqlEmbeds *embed.FS) {
    serverCmd := &cobra.Command{
        Use: "server",
        Short: "Start the server.",
        Run: func (cmd *cobra.Command, args []string) {
            pool, err := pgxpool.Connect(context.Background(), os.Getenv("PG_URL"))
            if err != nil {
                log.Fatal("PG_URL env var must be provided or invalid.")
            }
            defer pool.Close()

            ensureTablesCreated(pool, sqlEmbeds)
            ensureGrootUser(pool)
           
            fmt.Println("[racl] listening on:", listenAddr)
            log.Fatal(http.ListenAndServe(listenAddr, api.Router(pool)))
        },
    }

    serverCmd.
        Flags().
        StringVarP(&listenAddr, "listen-addr", "l", ":8080", "Listen address")

    rootCmd.AddCommand(serverCmd)
}

