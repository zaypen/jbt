package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/zaypen/jbt/core"
	"github.com/zaypen/jbt/util"
	"runtime"
)

func init() {
	rootCmd.AddCommand(&cobra.Command{
		Use:   "check",
		Short: "Check updates for installed products",
		Run: func(cmd *cobra.Command, args []string) {
			check()
		},
	})
}

func check() {
	u, err := core.New(runtime.GOOS)
	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
		return
	}
	installations := u.List(core.ProductCodes)
	updates := u.Check(installations)
	table := util.NewTable()
	table.Header("Code", "Product", "Installed", "Version", "Update", "Latest")
	for _, key := range core.ProductCodes {
		name := core.ProductNames[key]
		installation := installations[key]
		update, checked := updates[key]
		installed := util.If(installation.Installed, "yes", "no").(string)
		hasUpdate := util.If(checked, "yes", "no").(string)
		latestVersion := util.If(checked, update.Version, "").(string)
		table.Row(key, name, installed, installation.Version, hasUpdate, latestVersion)
	}
	table.End()
}
