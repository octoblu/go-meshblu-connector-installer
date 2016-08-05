package foreverizer

import "golang.org/x/sys/windows/registry"

func userLoginInstall(opts Options) error {
	key, err := registry.OpenKey(registry.CURRENT_USER, `Software\Microsoft\Windows\CurrentVersion\Run`, registry.WRITE)
	if err != nil {
		return err
	}
	key.SetStringValue(opts.ServiceName, opts.IgnitionPath)
	key.Close()
	return nil
}
