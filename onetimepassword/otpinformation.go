package onetimepassword

import "encoding/json"
import De "github.com/visionmedia/go-debug"

var debug = De.Debug("meshblu-connector-installer:onetimepassword:otpinformation")

// OTPInformation describes the information that a One Time Password can be exchanged for
type OTPInformation struct {
	UUID     string `json:"uuid"`
	Token    string `json:"token"`
	Metadata struct {
		GithubSlug      string `json:"githubSlug"`
		IgnitionVersion string `json:"ignitionVersion"`
		Connector       string `json:"connector"`
		Tag             string `json:"tag"`
		Meshblu         struct {
			Domain string `json:"domain,omitempty"`
		} `json:"meshblu,omitempty"`
	} `json:"metadata"`
}

// String returns the pretty json string representation of OTPInformation
func (info *OTPInformation) String() string {
	data, err := json.MarshalIndent(info, "", "  ")
	if err != nil {
		debug("Unexpected error: %v", err.Error())
		return "error"
	}

	return string(data)
}
