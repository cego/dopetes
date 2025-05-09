package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "dopetes",
	Short: "(d)(o)cker (p)ull (e)vents (t)o (e)lastic(s)earch",
	Long:  ``,
}

func Execute() {
	rootCmd.AddCommand(daemon)
	err := rootCmd.Execute()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
}
