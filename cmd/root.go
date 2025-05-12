package cmd

import (
	"github.com/cego/go-lib"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "dopetes",
	Short: "(d)(o)cker (p)ull (e)vents (t)o (e)lastic(s)earch",
}

func Execute() {
	logger := cego.NewLogger()
	rootCmd.AddCommand(InitDaemon())
	rootCmd.AddCommand(InitPublish())
	rootCmd.AddCommand(InitClear())
	err := rootCmd.Execute()
	if err != nil {
		logger.Error(err.Error())
		return
	}
}
