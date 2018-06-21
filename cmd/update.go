package cmd

import (
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(&cobra.Command{
		Use:   "update [product]",
		Short: "Update a product to latest version",
		Run: func(cmd *cobra.Command, args []string) {
			update()
		},
	})
}

func update() {
}
