package main

import (
    "embed"
    "github.com/bnert/racl/cmd"
)

//go:embed sql/schema/*.sql
var sqlEmbeds embed.FS

func main() {
    cmd.Execute(&sqlEmbeds)
}
