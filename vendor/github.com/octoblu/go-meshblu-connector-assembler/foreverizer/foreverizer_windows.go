package foreverizer

import (
	"fmt"
	"os/exec"

	"golang.org/x/sys/windows/registry"
)

func userLoginInstall(opts Options) error {
	ignitionPath := opts.IgnitionPath
	key, err := registry.OpenKey(registry.CURRENT_USER, `Software\Microsoft\Windows\CurrentVersion\Run`, registry.WRITE)
	if err != nil {
		return err
	}
	quotedIgnitionPath := fmt.Sprintf("\"%s\"", ignitionPath)
	debug("writing registry key...")
	key.SetStringValue(opts.ServiceName, quotedIgnitionPath)
	key.Close()

	cmd := exec.Command(ignitionPath)
	debug("starting...")
	err = cmd.Start()
	if err != nil {
		return err
	}
	return nil
}
