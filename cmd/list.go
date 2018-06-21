package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/zaypen/jbt/updater"
	"github.com/zaypen/jbt/util"
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
	u, err := updater.NewUpdater(runtime.GOOS)
	if err != nil {
		fmt.Printf("ERROR: %s", err.Error())
		return
	}
	installations := u.List()
	table := util.NewTable()
	table.Header("Code", "Product", "Installed", "Version")
	for _, key := range updater.ProductCodes {
		name := updater.ProductNames[key]
		installation := installations[key]
		table.Row(key, name, util.If(installation.Installed, "yes", "no").(string), installation.Version)
	}
	table.End()
}
