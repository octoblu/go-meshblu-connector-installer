package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"runtime/debug"
	"strings"

	"github.com/codegangsta/cli"
	"github.com/coreos/go-semver/semver"
	"github.com/fatih/color"
	"github.com/octoblu/go-meshblu-connector-assembler/assembler"
	"github.com/octoblu/go-meshblu-connector-dependency-manager/installer"
	"github.com/octoblu/go-meshblu-connector-installer/onetimepassword"
	"github.com/octoblu/go-meshblu-connector-installer/osruntime"
)

// NodeVersion is the version of node that will be installed
const NodeVersion = "v5.5.0"

// NPMVersion is the version of npm that will be installed (only applies for windows)
const NPMVersion = "v3.10.5"

const serviceTypeUserService = "UserService"
const serviceTypeService = "Service"

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
		cli.BoolFlag{
			Name:   "skip-one-time-password-expiration, s",
			EnvVar: "MESHBLU_CONNECTOR_INSTALLER_SKIP_ONE_TIME_PASSWORD_EXPIRATION",
			Usage:  "Skip the expiration of the one time password. This lets the installer run multiple times on the same password",
		},
		cli.StringFlag{
			Name:   "service-type, t",
			EnvVar: "MESHBLU_CONNECTOR_INSTALLER_SERVICE_TYPE",
			Usage:  "The type of install: Service, UserService, UserLogin. UserLogin is only available for Windows installs. Default: Service",
		},
		cli.StringFlag{
			Name:   "service-username, u",
			EnvVar: "MESHBLU_CONNECTOR_INSTALLER_SERVICE_USERNAME",
			Usage:  "Username to run the connector under",
		},
		cli.StringFlag{
			Name:   "service-password, p",
			EnvVar: "MESHBLU_CONNECTOR_INSTALLER_SERVICE_PASSWORD",
			Usage:  "Password for user account. Required for UserService",
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
	oneTimePassword, skipExpiration, serviceType, serviceUsername, servicePassword := getOpts(context)
	fmt.Println("Using One Time Password: ", oneTimePassword)

	if serviceType == serviceTypeUserService {
		if serviceUsername == "root" {
			color.Red("Refusing to install UserService as root, this is probably not what you want.")
			os.Exit(1)
		}
	}

	installNodeAndNPM(serviceType)
	connectorInfo, err := onetimepassword.GetOTPInformation(oneTimePassword)
	fatalIfError(err)

	UUID := connectorInfo.UUID
	Token := connectorInfo.Token

	ConnectorName := connectorInfo.Metadata.Connector
	GithubSlug := connectorInfo.Metadata.GithubSlug
	Tag := connectorInfo.Metadata.Tag
	IgnitionTag := connectorInfo.Metadata.IgnitionVersion

	Domain := connectorInfo.Metadata.Meshblu.Domain
	ResolveSRV := Domain != ""

	err = runAssembler(UUID, Token, ConnectorName, GithubSlug, Tag, IgnitionTag, serviceType, serviceUsername, servicePassword, Domain, ResolveSRV)
	fatalIfError(err)

	if skipExpiration {
		fmt.Fprintln(os.Stderr, "warning: skipping one-time-password expiration")
		os.Exit(0)
	}

	err = onetimepassword.Expire(oneTimePassword)
	fatalIfError(err)

	os.Exit(0)
}

func getOpts(context *cli.Context) (string, bool, string, string, string) {
	oneTimePassword := context.String("one-time-password")
	skipExpiration := context.Bool("skip-one-time-password-expiration")
	serviceType := context.String("service-type")
	serviceUsername := context.String("service-serviceUsername")
	servicePassword := context.String("service-password")

	if serviceType == "" {
		serviceType = "Service"
	}

	if oneTimePassword == "" {
		oneTimePassword = promptForOneTimePassword()
	}

	if oneTimePassword == "" {
		color.Red("meshblu-connector-installer needs a One Time Password to run")
		os.Exit(1)
	}

	if serviceType == serviceTypeUserService {
		if serviceUsername == "" {
			serviceUsername = os.Getenv("USER")
			if serviceUsername == "" {
				color.Red("meshblu-connector-installer needs a Username in UserService mode")
				os.Exit(1)
			}
		}
	}

	return oneTimePassword, skipExpiration, serviceType, serviceUsername, servicePassword
}

func installNodeAndNPM(serviceType string) {
	binPath, err := osruntime.UserBinPath(osruntime.New())
	if serviceType == serviceTypeService {
		binPath, err = osruntime.BinPath(osruntime.New())
	}
	fatalIfError(err)
	fatalIfError(installer.InstallNode(NodeVersion, binPath))
	fatalIfError(installer.InstallNPM(NPMVersion, binPath))
}

func promptForOneTimePassword() string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("One Time Password: ")
	text, err := reader.ReadString('\n')
	fatalIfError(err)
	return strings.TrimSpace(text)
}

func runAssembler(UUID, Token, ConnectorName, GithubSlug, Tag, IgnitionTag, ServiceType, ServiceUsername, ServicePassword, Domain string, ResolveSRV bool) error {
	options, err := assembler.NewOptions(assembler.OptionsOptions{
		ConnectorName:   ConnectorName,
		GithubSlug:      GithubSlug,
		Tag:             Tag,
		UUID:            UUID,
		Token:           Token,
		IgnitionTag:     IgnitionTag,
		ServiceType:     ServiceType,
		ServiceUsername: ServiceUsername,
		ServicePassword: ServicePassword,
		ResolveSRV:      ResolveSRV,
		Domain:          Domain,
	})

	if err != nil {
		return err
	}

	return assembler.Assemble(*options)
}

func version() string {
	version, err := semver.NewVersion(VERSION)
	if err != nil {
		errorMessage := fmt.Sprintf("Error with version number: %v", VERSION)
		log.Panicln(errorMessage, err.Error())
	}
	return version.String()
}
