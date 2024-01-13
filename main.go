package main

import (
	"github.com/artemys/pprof-visualizer/cmd/api"
	"github.com/spf13/cobra"
)

var rootCmd *cobra.Command

func init() {
	rootCmd = &cobra.Command{Use: "pprof-visualizer"}
	rootCmd.AddCommand(api.Start())
}

func main() {
	cobra.CheckErr(rootCmd.Execute())

}
