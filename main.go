package main

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"github.com/codegangsta/cli"
	"github.com/coreos/go-semver/semver"
	"github.com/fatih/color"
	"github.com/octoblu/go-meshblu-connector-installer/onetimepassword"
	De "github.com/visionmedia/go-debug"
)

var debug = De.Debug("meshblu-connector-installer:main")

func main() {
	app := cli.NewApp()
	app.Name = "meshblu-connector-installer"
	app.Version = version()
	app.Action = run
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "one-time-password, o",
			EnvVar: "MESHBLU_CONNECTOR_INSTALLER_ONE_TIME_PASSWORD",
			Usage:  "The one time password provided by the Connector Factory",
		},
	}
	app.Run(os.Args)
}

func fatalIfError(err error) {
	if err == nil {
		return
	}
	log.Fatalln("Fatal Error", err.Error())
}

func run(context *cli.Context) {
	oneTimePassword := getOpts(context)
	fmt.Println("Using One Time Password: ", oneTimePassword)

	connectorInfo, err := onetimepassword.GetOTPInformation(oneTimePassword)
	fatalIfError(err)
	fmt.Println("Got info: ", connectorInfo)

	err := dependencies.Install()
	fatalIfError(err)
}

func getOpts(context *cli.Context) string {
	oneTimePassword := context.String("one-time-password")

	if oneTimePassword == "" {
		oneTimePassword = promptForOneTimePassword()
	}

	if oneTimePassword == "" {
		color.Red("meshblu-connector-installer needs a One Time Password to run")
		os.Exit(1)
	}

	return oneTimePassword
}

func promptForOneTimePassword() string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("One Time Password: ")
	text, err := reader.ReadString('\n')
	fatalIfError(err)
	return text
}

func version() string {
	version, err := semver.NewVersion(VERSION)
	if err != nil {
		errorMessage := fmt.Sprintf("Error with version number: %v", VERSION)
		log.Panicln(errorMessage, err.Error())
	}
	return version.String()
}
