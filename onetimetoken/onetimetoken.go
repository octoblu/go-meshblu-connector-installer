package onetimetoken

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
)

// OTP is an interface for retrieving information about One Time Tokens and to expire them
type OTP interface {
	ExchangeForInformation() OTPInformation
}

type httpOTP struct {
	oneTimePassword string
	baseURL         *url.URL
}

// OTPInformation describes the information that a One Time Password can be exchanged for
type OTPInformation struct {
	UUID, Token string
}

// New constructs a new OTP instance
func New(oneTimePassword string) OTP {
	otp, err := NewWithURLOverride(oneTimePassword, "https://meshblu-otp.octoblu.com/")
	if err != nil {
		log.Fatalln("This URL should never be invalid, but it is: ", err.Error())
	}
	return otp
}

// NewWithURLOverride constructs a new OTP instance with a specific URL
func NewWithURLOverride(oneTimePassword, urlStr string) (OTP, error) {
	baseURL, err := url.ParseRequestURI(urlStr)

	if err != nil {
		return nil, err
	}

	return &httpOTP{oneTimePassword, baseURL}, nil
}

// ExchangeForInformation exchanges a one time token for information, including
// the connector type and Meshblu credentials
func (otp *httpOTP) ExchangeForInformation() OTPInformation {
	retrievalURL := *otp.baseURL
	retrievalURL.Path = fmt.Sprintf("/retrieve/%v", otp.oneTimePassword)

	response, _ := http.Get(retrievalURL.String())
	return parseRetrievalResponse(response)
}

func parseRetrievalResponse(response *http.Response) OTPInformation {
	return OTPInformation{}
}

// uuid: 'c7097087-bed4-4272-8692-3b07277ec281',
// token: 'a7702204e157e51fd63c924a7b77db00121a0d5d',
// metadata: {
//   githubSlug: 'octoblu/meshblu-connector-say-hello',
//   connectorAssemblerVersion: 'v13.0.0',
//   dependencyManagerVersion: 'v3.0.2',
//   ignitionVersion: 'v6.0.0',
//   connector: 'say-hello',
//   tag: 'v6.0.0',
// },
