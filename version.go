package main

import (
	"fmt"
	"runtime"

	"github.com/spf13/cobra"
)

var version = "compiled manually"

var cmdVersion = &cobra.Command{
	Use:   "version",
	Short: "Display version information",
	Run: func(*cobra.Command, []string) {
		fmt.Printf("elmo %s\ncompiled with %v on %v\n",
			version, runtime.Version(), runtime.GOOS)
	},
}
