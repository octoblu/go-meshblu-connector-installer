package main

import (
	"fmt"
	"log"
	"os"

	"github.com/codegangsta/cli"
	"github.com/coreos/go-semver/semver"
	"github.com/fatih/color"
	"github.com/octoblu/go-meshblu-connector-dependency-manager/installer"
	De "github.com/tj/go-debug"
)

var debug = De.Debug("meshblu-connector-dependency-manager:main")

func main() {
	app := cli.NewApp()
	app.Name = "meshblu-connector-dependency-manager"
	app.Version = version()
	app.Action = run
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "tag",
			EnvVar: "MESHBLU_CONNECTOR_DEPENDENCY_TAG",
			Usage:  "Version of the dependency you want. Values: 'node', 'npm', or 'nssm'.",
		},
		cli.StringFlag{
			Name:   "type",
			EnvVar: "MESHBLU_CONNECTOR_DEPENDENCY_TYPE",
			Usage:  "Type of dependency you want",
		},
	}
	app.Run(os.Args)
}

func run(context *cli.Context) {
	tag, depType := getOpts(context)

	installerClient := installer.New()
	err := installerClient.Do(depType, tag)
	fatalIfError("error installing dep:", err)
}

func getOpts(context *cli.Context) (string, string) {
	tag := context.String("tag")
	depType := context.String("type")

	if tag == "" || depType == "" {
		cli.ShowAppHelp(context)

		if tag == "" {
			color.Red("  Missing required flag --tag or MESHBLU_CONNECTOR_DEPENDENCY_MANAGER_TAG")
		}

		if depType == "" {
			color.Red("  Missing required flag --type or MESHBLU_CONNECTOR_DEPENDENCY_MANAGER_TYPE")
		}

		os.Exit(1)
	}

	return tag, depType
}

func version() string {
	version, err := semver.NewVersion(VERSION)
	if err != nil {
		errorMessage := fmt.Sprintf("Error with version number: %v", VERSION)
		log.Panicln(errorMessage, err.Error())
	}
	return version.String()
}

func fatalIfError(msg string, err error) {
	if err == nil {
		return
	}

	log.Fatalln(msg, err.Error())
}
