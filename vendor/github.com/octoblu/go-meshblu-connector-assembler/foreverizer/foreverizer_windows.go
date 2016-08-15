package foreverizer

import (
	"fmt"

	"golang.org/x/sys/windows/registry"
)

func userLoginInstall(opts Options) error {
	key, err := registry.OpenKey(registry.CURRENT_USER, `Software\Microsoft\Windows\CurrentVersion\Run`, registry.WRITE)
	if err != nil {
		return err
	}
	ignitionPath := fmt.Sprintf("\"%s\"", opts.IgnitionPath)
	key.SetStringValue(opts.ServiceName, ignitionPath)
	key.Close()
	return nil
}
