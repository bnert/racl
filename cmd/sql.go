package cmd

import (
    "fmt"
    "embed"
    "path"
    "github.com/spf13/cobra"
)

func registerSqlCmd(sqlEmbeds *embed.FS) {
    sqlCmd := &cobra.Command{
        Use: "sql",
        Short: "List sql defintions",
        Run: func(cmd *cobra.Command, args []string) {
            basePath := "sql/schema"

            dir, err := sqlEmbeds.ReadDir(basePath)
            if err != nil {
                fmt.Println("Cannot open file", basePath)
                return
            }

            for _, v := range dir {
                p := path.Join(basePath, v.Name())
                f, err := sqlEmbeds.ReadFile(p)
                if err != nil {
                    fmt.Println("Cannot open file:", p)
                }
                fmt.Println("--", v.Name())
                fmt.Println(string(f))
            }
        },
    }

    rootCmd.AddCommand(sqlCmd)
}

