package dependencies

// Runtime is used to provide the GOOS and GOARCH variables to the DownloadWithURLAndRuntime
// Generally, this should not be instantiated by hand. Instead, call Download() and it'll take
// care of this part for you
type Runtime struct {
	GOOS, GOARCH string
}
