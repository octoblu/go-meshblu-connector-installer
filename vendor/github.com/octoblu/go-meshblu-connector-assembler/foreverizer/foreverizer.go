package foreverizer

import (
	"fmt"

	"github.com/kardianos/service"
	De "github.com/tj/go-debug"
)

var debug = De.Debug("meshblu-connector-assembler:foreverizer")

// Options device the options to be passed
// to construct a new Foreverizer
type Options struct {
	ServiceName     string
	ServiceType     string
	ServiceUsername string
	ServicePassword string
	DisplayName     string
	Description     string
	IgnitionPath    string
}

// Foreverize registers the service with the OS
func Foreverize(opts Options) error {
	err := validateOptions(opts)
	if err != nil {
		return err
	}

	debug("foreverizing...")
	fmt.Println("ServiceName", opts.ServiceName)
	if opts.ServiceType == "UserLogin" {
		return userLoginInstall(opts)
	}

	return serviceInstall(opts)
}

func serviceInstall(opts Options) error {
	userService := false
	if opts.ServiceType == "UserService" {
		userService = true
	}
	svcConfig := &service.Config{
		Name:        opts.ServiceName,
		DisplayName: opts.DisplayName,
		Description: opts.Description,
		Executable:  opts.IgnitionPath,
		Dependencies: []string{ "LanmanServer" },
		Option: service.KeyValue{
			"UserService": userService,
			"KeepAlive":   true,
			"UserName":    opts.ServiceUsername,
			"Password":    opts.ServicePassword,
		},
	}

	prg := &Program{}
	connectorService, err := service.New(prg, svcConfig)
	if err != nil {
		return err
	}

	debug("maybe stop and removing service...")
	err = connectorService.Uninstall()
	if err != nil {
		debug("Error on uninstall, it probably wasn't installed in the first place: %v", err.Error())
	}

	debug("installing service...")
	err = connectorService.Install()
	if err != nil {
		return err
	}
	debug("starting...")
	err = connectorService.Start()
	if err != nil {
		return err
	}
	return nil
}

func validateOptions(opts Options) error {
	if opts.Description == "" {
		return fmt.Errorf("Missing required param: Description")
	}
	if opts.DisplayName == "" {
		return fmt.Errorf("Missing required param: DisplayName")
	}
	if opts.IgnitionPath == "" {
		return fmt.Errorf("Missing required param: IgnitionPath")
	}
	if opts.ServiceName == "" {
		return fmt.Errorf("Missing required param: ServiceName")
	}
	if opts.ServiceType == "" {
		return fmt.Errorf("Missing required param: ServiceType")
	}
	if opts.ServiceType == "UserService" {
		if opts.ServiceUsername == "" {
			return fmt.Errorf("Missing required param: ServiceUsername")
		}
	}

	return nil
}
