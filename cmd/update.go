package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/zaypen/jbt/core"
	"github.com/zaypen/jbt/util"
	"os"
	"runtime"
	"strings"
)

func init() {
	cmd := &cobra.Command{
		Use:   "update [product]",
		Short: "Update a product to latest version",
		Run: func(cmd *cobra.Command, args []string) {
			update(args)
		},
	}
	rootCmd.AddCommand(cmd)
}

func parseCodes(codes []string) []string {
	formattedCodes := codes[:0]
	for _, code := range codes {
		formattedCodes = append(formattedCodes, strings.ToUpper(code))
	}
	for i, code := range formattedCodes {
		if !util.In(code, core.ProductCodes) {
			fmt.Printf("Unknown product code: %s\n", codes[i])
			os.Exit(1)
		}
	}
	if len(codes) == 0 {
		codes = core.ProductCodes
	}
	return codes
}

func update(codes []string) {
	codes = parseCodes(codes)
	u, err := core.New(runtime.GOOS)
	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
		os.Exit(255)
	}
	u.Update(u.Check(u.List(codes)))
}
