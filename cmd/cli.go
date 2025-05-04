package cmd

import (
	"fmt"

	"github.com/neilfarmer/go-git/internal/config"
	"github.com/neilfarmer/go-git/internal/github"
	"github.com/neilfarmer/go-git/internal/gitlab"
	"github.com/spf13/cobra"
)


var syncCmd = &cobra.Command{
	Use:   "sync",
	Short: "Sync repositories from SCM",
	Long:  `Sync repositories from SCM`,
	Run: func(cmd *cobra.Command, args []string) {
		config, err := config.ReadConfig()
		if err != nil {
			fmt.Println("Error reading config:", err)
			return
		}
		if config.SCM_Name == "github" {
			fmt.Println("GitHub configuration detected.")
			github.GetRepos(config)
		}
		if config.SCM_Name == "gitlab" {
			fmt.Println("Gitlab configuration detected.")
			gitlab.GetRepos(config)
		}
	},
}
var graphCmd = &cobra.Command{
	Use:   "graph",
	Short: "Graph repositories from SCM",
	Long:  `Graph repositories from SCM`,
	Run: func(cmd *cobra.Command, args []string) {
		config, err := config.ReadConfig()
		if err != nil {
			fmt.Println("Error reading config:", err)
			return
		}
		if config.SCM_Name == "gitlab" {
			fmt.Println("Gitlab configuration detected.")
			gitlab.GraphRepos(config)
			// gitlab.GetRepos(config)
		}
	},
}
func init() {
	rootCmd.AddCommand(syncCmd)
	rootCmd.AddCommand(graphCmd)
	rootCmd.PersistentFlags().StringVarP(&verboseLevel, "verbose", "v", "info", "Log level: debug, info, warn, error")
}