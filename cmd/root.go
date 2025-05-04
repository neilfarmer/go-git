package cmd

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/neilfarmer/go-git/internal"
	"github.com/spf13/cobra"
)

var verboseLevel string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "go-git",
	Short: "Go-git is a command line tool to clone and manage repositories",
	Long: `Go-git is a command line tool to clone and manage repositories from different source control management systems (SCM) like GitHub, GitLab, etc.`,

	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		var level slog.Level
		switch verboseLevel {
		case "debug":
			level = slog.LevelDebug
		case "info":
			level = slog.LevelInfo
		case "warn":
			level = slog.LevelWarn
		case "error":
			level = slog.LevelError
		default:
			level = slog.LevelInfo
		}
		logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
			Level: level,
		}))
		slog.SetDefault(logger)
	},
}

func Execute() {
	// Couldnt figure out how to do this the cobra way
	for _, arg := range os.Args[1:] {
		if arg == "-V" || arg == "--version" {
			fmt.Println("go-git version:", internal.Version)
			os.Exit(0)
		}
	}

	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	rootCmd.PersistentFlags().BoolP("version", "V", false, "Print version and exit")
}
