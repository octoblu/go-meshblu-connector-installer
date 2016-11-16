package main

import (
	"fmt"
	"log"
	"os"

	"github.com/codegangsta/cli"
	"github.com/coreos/go-semver/semver"
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
			Usage:  "Version of the dependency you want",
		},
		cli.StringFlag{
			Name:   "type",
			EnvVar: "MESHBLU_CONNECTOR_DEPENDENCY_TYPE",
			Usage:  "Type of dependency you want. Values: 'node', 'npm', or 'nssm'.",
		},
	}
	app.Run(os.Args)
}

func run(context *cli.Context) {
	log.Fatalln("DEPRACTION NOTICE", fmt.Errorf("Must be used as a library").Error())
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
