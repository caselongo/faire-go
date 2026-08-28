package faire_go

import (
	"bytes"
	"encoding/json"
	"fmt"
	errortools "github.com/leapforce-libraries/go_errortools"
	go_http "github.com/leapforce-libraries/go_http"
	oauth2 "github.com/leapforce-libraries/go_oauth2"
	"github.com/leapforce-libraries/go_oauth2/tokenfixed"
	"net/http"
	"strings"
)

const (
	apiUrl             string = "https://www.faire.com/external-api/v2"
	defaultRedirectUrl string = "http://localhost:8080/oauth/redirect"
	authUrl            string = "https://faire.com/oauth2/authorize"
	tokenUrl           string = "https://www.faire.com/api/external-api-oauth2/token"
	tokenHttpMethod    string = http.MethodPost
)

type Service struct {
	applicationId     string
	applicationSecret string
	appCredentials    string
	accessToken       string
	oAuth2Service     *oauth2.Service
	redirectUrl       *string
	errorResponse     *ErrorResponse
}

type ServiceConfig struct {
	ApplicationId     string
	ApplicationSecret string
	AppCredentials    string
	AccessToken       string
	RedirectUrl       *string
}

type TokenRequest struct {
	ApplicationToken  string   `json:"application_token"`
	ApplicationSecret string   `json:"application_secret"`
	RedirectUrl       string   `json:"redirect_url"`
	Scope             []string `json:"scope"`
	GrantType         string   `json:"grant_type"`
	AuthorizationCode string   `json:"authorization_code"`
}

func (service *Service) getTokenRequest(r *http.Request) (*http.Request, *errortools.Error) {
	err := r.ParseForm()
	if err != nil {
		return nil, errortools.ErrorMessage(err)
	}
	authorizationCode := r.FormValue("authorizationCode")

	reqBody := TokenRequest{
		ApplicationToken:  service.applicationId,
		ApplicationSecret: service.applicationSecret,
		RedirectUrl:       *service.redirectUrl,
		Scope: []string{
			"READ_BRAND",
			"READ_PRODUCTS",
			"WRITE_PRODUCTS",
			"READ_RETAILER",
			"READ_ORDERS",
		},
		GrantType:         "AUTHORIZATION_CODE",
		AuthorizationCode: authorizationCode,
	}

	body, err := json.Marshal(reqBody)
	if err != nil {
		panic(err)
	}

	req, err := http.NewRequest(
		http.MethodPost,
		tokenUrl,
		bytes.NewBuffer(body),
	)
	if err != nil {
		return nil, errortools.ErrorMessage(err)
	}

	req.Header.Set("Content-Type", "application/json")

	return req, nil
}

func NewService(cfg *ServiceConfig) (*Service, *errortools.Error) {
	if cfg == nil {
		return nil, errortools.ErrorMessage("ServiceConfig must not be a nil pointer")
	}

	if cfg.ApplicationId == "" {
		return nil, errortools.ErrorMessage("ApplicationId not provided")
	}

	redirectUrl := defaultRedirectUrl
	if cfg.RedirectUrl != nil {
		redirectUrl = *cfg.RedirectUrl
	}

	var service = Service{
		applicationId:     cfg.ApplicationId,
		applicationSecret: cfg.ApplicationSecret,
		appCredentials:    cfg.AppCredentials,
		accessToken:       cfg.AccessToken,
		redirectUrl:       cfg.RedirectUrl,
	}

	tokenSource, e := tokenfixed.NewTokenFixed(cfg.AccessToken)
	if e != nil {
		panic(e)
	}

	var getTokenRequestFunc = service.getTokenRequest

	oauth2ServiceConfig := oauth2.ServiceConfig{
		ClientId:                cfg.ApplicationId,
		ClientSecret:            cfg.ApplicationSecret,
		RedirectUrl:             redirectUrl,
		AuthUrl:                 authUrl,
		TokenUrl:                tokenUrl,
		TokenHttpMethod:         tokenHttpMethod,
		TokenSource:             tokenSource,
		GetTokenFromRequestFunc: &getTokenRequestFunc,
	}
	oauth2Service, e := oauth2.NewService(&oauth2ServiceConfig)
	if e != nil {
		return nil, e
	}

	service.oAuth2Service = oauth2Service

	return &service, nil
}

func (service *Service) httpRequest(requestConfig *go_http.RequestConfig) (*http.Request, *http.Response, *errortools.Error) {
	// add error model
	service.errorResponse = &ErrorResponse{}
	requestConfig.ErrorModel = service.errorResponse

	// add authentication headers
	header := http.Header{}
	header.Set("X-FAIRE-APP-CREDENTIALS", service.appCredentials)
	header.Set("X-FAIRE-OAUTH-ACCESS-TOKEN", service.accessToken)
	(*requestConfig).NonDefaultHeaders = &header

	req, res, e := service.oAuth2Service.HttpRequest(requestConfig)
	if e != nil {

		if service.errorResponse != nil {
			b, err := json.Marshal(service.errorResponse)
			if err == nil {
				e.SetMessage(string(b))
			}
		}
	}

	return req, res, e
}

func (service *Service) AuthorizeUrl(scopes []string, state string) string {
	if service.redirectUrl == nil {
		return ""
	}
	return fmt.Sprintf("%s?applicationId=%s&redirectUrl=%s&state=%s&scope=%s", authUrl, service.applicationId, *service.redirectUrl, state, strings.Join(scopes, "&scope="))
}

func (service *Service) GetTokenFromCode(r *http.Request) *errortools.Error {
	return service.oAuth2Service.GetTokenFromCode(r, nil)
}

func (service *Service) url(path string) string {
	return fmt.Sprintf("%s/%s", apiUrl, path)
}

func (service *Service) ErrorResponse() *ErrorResponse {

	return service.errorResponse
}

type Amount struct {
	AmountMinor float64 `json:"amount_minor"`
	Currency    string  `json:"currency"`
}

type Option struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type GeoConstraint struct {
	CountryGroup string `json:"country_group"`
}
