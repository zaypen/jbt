package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/zaypen/jbt/core"
	"github.com/zaypen/jbt/util"
	"os"
	"runtime"
)

func init() {
	rootCmd.AddCommand(&cobra.Command{
		Use:   "list",
		Short: "Print installed products with status",
		Run: func(cmd *cobra.Command, args []string) {
			list()
		},
	})
}

func list() {
	u, err := core.New(runtime.GOOS)
	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
		os.Exit(255)
	}
	installations := u.List(core.ProductCodes)
	table := util.NewTable()
	table.Header("Code", "Product", "Installed", "Version")
	for _, key := range core.ProductCodes {
		name := core.ProductNames[key]
		installation := installations[key]
		table.Row(key, name, util.If(installation.Installed, "yes", "no").(string), installation.Version)
	}
	table.End()
}
