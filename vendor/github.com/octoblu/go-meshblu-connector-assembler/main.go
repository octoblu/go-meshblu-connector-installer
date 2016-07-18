package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"runtime"

	"github.com/codegangsta/cli"
	"github.com/coreos/go-semver/semver"
	"github.com/octoblu/go-meshblu-connector-assembler/assembler"
	De "github.com/tj/go-debug"
)

var debug = De.Debug("meshblu-connector-assembler:main")

func main() {
	app := cli.NewApp()
	app.Name = "meshblu-connector-assembler"
	app.Version = version()
	app.Action = run
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "connector, c",
			EnvVar: "MESHBLU_CONNECTOR_ASSEMBLER_CONNECTOR",
			Usage:  "Connector name",
		},
		cli.BoolFlag{
			Name:  "debug, d",
			Usage: "Debug mode, will prompt for user to continue on Windows",
		},
		cli.StringFlag{
			Name:   "github-slug, g",
			EnvVar: "MESHBLU_CONNECTOR_ASSEMBLER_GITHUB_SLUG",
			Usage:  "Github Slug",
		},
		cli.StringFlag{
			Name:   "tag, T",
			EnvVar: "MESHBLU_CONNECTOR_ASSEMBLER_TAG",
			Usage:  "Tag or Version",
		},
		cli.StringFlag{
			Name:   "ignition, i",
			EnvVar: "MESHBLU_CONNECTOR_ASSEMBLER_IGNITION_TAG",
			Usage:  "Ignition Tag",
		},
		cli.StringFlag{
			Name:   "uuid, u",
			EnvVar: "MESHBLU_CONNECTOR_ASSEMBLER_UUID",
			Usage:  "Meshblu device uuid",
		},
		cli.StringFlag{
			Name:   "token, t",
			EnvVar: "MESHBLU_CONNECTOR_ASSEMBLER_TOKEN",
			Usage:  "Meshblu device token",
		},
	}
	app.Run(os.Args)
}

func run(context *cli.Context) {
	assemblerOptions, debugEnabled := getOpts(context)
	err := assembler.Assemble(*assemblerOptions)
	fatalIfError(debugEnabled, "Error assembling", err)

	debug("done installing")
	if debugEnabled {
		tellWindowsToWait()
	}
}

func fatalIfError(windowsShouldWait bool, msg string, err error) {
	if err == nil {
		return
	}

	log.Println(msg, err.Error())
	if windowsShouldWait {
		tellWindowsToWait()
	}
	log.Fatalln("Exiting...")
}

func getOpts(context *cli.Context) (*assembler.Options, bool) {
	assemblerOptions, err := assembler.NewOptions(assembler.OptionsOptions{
		ConnectorName: context.String("connector"),
		GithubSlug:    context.String("github-slug"),
		Tag:           context.String("tag"),
		UUID:          context.String("uuid"),
		Token:         context.String("token"),
		IgnitionTag:   context.String("ignition"),
	})
	fatalIfError(true, "Error populating default options", err)
	return assemblerOptions, context.Bool("Debug")
}

func tellWindowsToWait() {
	if runtime.GOOS != "windows" {
		return
	}

	fmt.Println("Press any key to continue >>>")
	bufio.NewReader(os.Stdin).ReadBytes('\n')
}

func version() string {
	version, err := semver.NewVersion(VERSION)
	if err != nil {
		errorMessage := fmt.Sprintf("Error with version number: %v", VERSION)
		log.Panicln(errorMessage, err.Error())
	}
	return version.String()
}
