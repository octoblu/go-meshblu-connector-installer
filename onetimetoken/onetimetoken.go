package onetimetoken

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)

// OTP is an interface for retrieving information about One Time Tokens and to expire them
type OTP interface {
	// ExchangeForInformation exchanges a one time token for information, including
	// the connector type and Meshblu credentials
	ExchangeForInformation() (*OTPInformation, error)
}

type httpOTP struct {
	oneTimePassword string
	baseURL         *url.URL
}

// OTPInformation describes the information that a One Time Password can be exchanged for
type OTPInformation struct {
	UUID  string `json:"uuid"`
	Token string `json:"token"`
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
func (otp *httpOTP) ExchangeForInformation() (*OTPInformation, error) {
	retrievalURL := *otp.baseURL
	retrievalURL.Path = fmt.Sprintf("/retrieve/%v", otp.oneTimePassword)

	response, err := http.Get(retrievalURL.String())
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	if response.StatusCode != 200 {
		return nil, fmt.Errorf("Received non 200: %v, %v", response.StatusCode, body)
	}

	return parseRetrievalResponse(body)
}

func parseRetrievalResponse(body []byte) (*OTPInformation, error) {
	info := OTPInformation{}

	err := json.Unmarshal(body, &info)
	if err != nil {
		return nil, err
	}

	return &info, nil
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
