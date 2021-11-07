package cmd

import (
    "embed"
    "github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
    Use: "racl",
    Short: "rest access control lists",
    Run: func(c *cobra.Command, args []string) {
        // Stuff?
    },
}

func Execute(sqlEmbeds *embed.FS) {
    registerServerCmd(sqlEmbeds)
    registerSqlCmd(sqlEmbeds)
    rootCmd.Execute()
}
