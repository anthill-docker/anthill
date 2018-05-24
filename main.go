/*
Package main is the entry point of the application.

TODO
*/
package main

import (
	"fmt"
	"os"
	"time"

	"github.com/aenthill/aenthill/commands"

	"github.com/apex/log"
	"github.com/apex/log/handlers/cli"
)

/*
version will be set by GoReleaser.
It will be the current Git tag (with v prefix stripped) or
the name of the snapshot if you're using the --snapshot flag.
*/
var version = "master"

func init() {
	log.SetHandler(cli.Default)
	commands.RootCmd.Version = version
}

func main() {
	fmt.Println()

	start := time.Now()
	if err := commands.RootCmd.Execute(); err != nil {
		log.WithError(err).Errorf("aenthill command failed after %0.2fs", time.Since(start).Seconds())
		fmt.Println()
		os.Exit(1)
	}

	hasCommand := false
	for _, arg := range os.Args {
		if arg == commands.InitCmd.Use || arg == commands.AddCmd.Use || arg == commands.RemoveCmd.Use {
			hasCommand = true
		}
	}

	hasHelpFlag := false
	for _, arg := range os.Args {
		if arg == "-h" || arg == "--help" {
			hasHelpFlag = true
		}
	}

	if len(os.Args) > 1 && !hasHelpFlag && hasCommand {
		log.Infof("aenthill command finished after %0.2fs", time.Since(start).Seconds())
	}
	fmt.Println()
}