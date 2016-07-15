package installer

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// FinalDependencyFileName gets the final dependency filename
func FinalDependencyFileName(depType, tag string) string {
	if depType == NodeType {
		return "node.exe"
	}
	if depType == NSSMType {
		return "nssm.exe"
	}
	if depType == NPMType {
		return "npm"
	}
	return ""
}

// GetResourceURI defines the uri to download
func GetResourceURI(depType, tag string) string {
	if depType == NodeType {
		return getNodeURI(tag)
	}
	if depType == NSSMType {
		return getNSSMURI(tag)
	}
	if depType == NPMType {
		return getNPMURI(tag)
	}
	return ""
}

func getNSSMURI(tag string) string {
	return fmt.Sprintf("http://nssm.cc/release/nssm-%v.zip", tag)
}

func getNodeURI(tag string) string {
	return strings.Replace("https://nodejs.org/dist/:tag:/win-x64/node.exe", ":tag:", tag, -1)
}

func getNPMURI(tag string) string {
	return fmt.Sprintf("https://github.com/npm/npm/archive/%s.zip", tag)
}

// GetBinPath defines the target location
func GetBinPath() string {
	return filepath.Join(os.Getenv("LOCALAPPDATA"), "MeshbluConnectors", "bin")
}

// ExtractBin allows you too extract the bin from the download
func ExtractBin(depType, target, tag string) error {
	if depType == NodeType {
		return nil
	}
	if depType == NSSMType {
		return ExtractNSSM(target, tag)
	}
	if depType == NPMType {
		return ExtractNPM(target, tag)
	}
	return fmt.Errorf("Unsupported platform")
}

// ExtractNSSM extracts the unzipped nssm directory
func ExtractNSSM(target, tag string) error {
	folderName := fmt.Sprintf("nssm-%s", tag)
	nssmPath := filepath.Join(target, folderName, "win64", "nssm.exe")
	nssmNewPath := filepath.Join(target, "nssm.exe")
	err := CopyFile(nssmPath, nssmNewPath)
	if err != nil {
		return err
	}
	return nil
}

// ExtractNPM extracts the unzipped nssm directory
func ExtractNPM(target, tag string) error {
	folderName := fmt.Sprintf("npm-%s", strings.Replace(tag, "v", "", -1))
	npmPath := filepath.Join(target, folderName)
	nodeModulesPath := filepath.Join(target, "node_modules")
	err := os.MkdirAll(nodeModulesPath, os.ModePerm)
	if err != nil {
		return err
	}

	err = os.Rename(filepath.Join(npmPath, "bin", "npm"), filepath.Join(target, "npm"))
	if err != nil {
		return fmt.Errorf("Error renaming npm %v", err.Error())
	}

	err = os.Rename(filepath.Join(npmPath, "bin", "npm.cmd"), filepath.Join(target, "npm.cmd"))
	if err != nil {
		return fmt.Errorf("Error renaming npm.cmd %v", err.Error())
	}

	npmNewPath := filepath.Join(nodeModulesPath, "npm")
	os.RemoveAll(npmNewPath)

	err = os.Rename(npmPath, npmNewPath)
	if err != nil {
		return fmt.Errorf("Error renaming npm node_modules %v", err.Error())
	}
	return nil
}
