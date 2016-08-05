package foreverizer

import (
	"fmt"
	"runtime"
)

func userLoginInstall(opts Options) error {
	return fmt.Errorf("UserLogin is not available on this platform: %s", runtime.GOOS)
}
