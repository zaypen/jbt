package cmd

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/zaypen/jbt/util"
	"os"
	"path/filepath"
)

var verbose bool

var rootCmd = &cobra.Command{
	Use:     filepath.Base(os.Args[0]),
	Short:   "JetBrains tools",
	Long:    "JBT is a toolbox to manage your JetBrains products",
	Version: "0.1.0",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		logrus.SetFormatter(&logrus.TextFormatter{})
		logrus.SetLevel(util.If(verbose, logrus.DebugLevel, logrus.InfoLevel).(logrus.Level))
	},
}

func init() {
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Make the operation more talkative")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
