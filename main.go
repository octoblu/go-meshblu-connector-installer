package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"runtime"
	"runtime/debug"

	"github.com/codegangsta/cli"
	"github.com/coreos/go-semver/semver"
	"github.com/fatih/color"
	"github.com/octoblu/go-meshblu-connector-assembler/assembler"
	"github.com/octoblu/go-meshblu-connector-dependency-manager/installer"
	"github.com/octoblu/go-meshblu-connector-installer/onetimepassword"
)

// NodeVersion is the version of node that will be installed
const NodeVersion = "v5.5.0"

// NPMVersion is the version of npm that will be installed (only applies for windows)
const NPMVersion = "v3.10.5"

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
	debug.PrintStack()
	log.Fatalln("Fatal Error", err.Error())
}

func run(context *cli.Context) {
	oneTimePassword := getOpts(context)
	fmt.Println("Using One Time Password: ", oneTimePassword)

	installNodeAndNPM()
	connectorInfo, err := onetimepassword.GetOTPInformation(oneTimePassword)
	fatalIfError(err)

	UUID := connectorInfo.UUID
	Token := connectorInfo.Token

	ConnectorName := connectorInfo.Metadata.Connector
	GithubSlug := connectorInfo.Metadata.GithubSlug
	Tag := connectorInfo.Metadata.Tag
	IgnitionTag := connectorInfo.Metadata.IgnitionVersion

	runAssembler(UUID, Token, ConnectorName, GithubSlug, Tag, IgnitionTag)

	fmt.Println("Got info: ", connectorInfo.String())

	os.Exit(0)
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

func installNodeAndNPM() {
	fatalIfError(installer.Install("node", "v5.5.0"))

	if runtime.GOOS == "windows" {
		fatalIfError(installer.Install("npm", "v3.10.5"))
	}
}

func promptForOneTimePassword() string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("One Time Password: ")
	text, err := reader.ReadString('\n')
	fatalIfError(err)
	return text
}

func runAssembler(UUID, Token, ConnectorName, GithubSlug, Tag, IgnitionTag string) {
	options, err := assembler.NewOptions(assembler.OptionsOptions{
		ConnectorName: ConnectorName,
		GithubSlug:    GithubSlug,
		Tag:           Tag,
		UUID:          UUID,
		Token:         Token,
		IgnitionTag:   IgnitionTag,
	})
	fatalIfError(err)
	err = assembler.Assemble(*options)
	fatalIfError(err)
}

func version() string {
	version, err := semver.NewVersion(VERSION)
	if err != nil {
		errorMessage := fmt.Sprintf("Error with version number: %v", VERSION)
		log.Panicln(errorMessage, err.Error())
	}
	return version.String()
}
