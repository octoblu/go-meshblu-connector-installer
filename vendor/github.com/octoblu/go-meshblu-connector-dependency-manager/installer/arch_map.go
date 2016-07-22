package installer

// ArchMap defines the GOARCH to NodeArch mapping
var ArchMap = map[string]string{
	"386":   "x86",
	"amd64": "x64",
	"arm":   "armv7l",
}
