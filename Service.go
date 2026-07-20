package faire_go

import (
	"fmt"
	errortools "github.com/leapforce-libraries/go_errortools"
	go_http "github.com/leapforce-libraries/go_http"
	oauth2 "github.com/leapforce-libraries/go_oauth2"
	"github.com/leapforce-libraries/go_oauth2/tokensource"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

const (
	apiUrl             string = "https://www.faire.com/external-api/v2"
	defaultRedirectUrl string = "http://localhost:8080/oauth/redirect"
	authUrl            string = "https://faire.com/oauth2/authorize"
	tokenUrl           string = "https://www.faire.com/api/external-api-oauth2/token"
	tokenHttpMethod    string = http.MethodPost
)

type Service struct {
	clientId      string
	clientSecret  string
	oAuth2Service *oauth2.Service
	redirectUrl   *string
	errorResponse *ErrorResponse
}

type ServiceConfig struct {
	ClientId     string
	ClientSecret string
	TokenSource  tokensource.TokenSource
	RedirectUrl  *string
}

func (service *Service) getTokenRequest(r *http.Request) (*http.Request, *errortools.Error) {
	err := r.ParseForm()
	if err != nil {
		return nil, errortools.ErrorMessage(err)
	}
	code := r.FormValue("code")

	data := url.Values{}
	data.Set("application_token", service.clientId)
	data.Set("application_secret", service.clientSecret)
	data.Set("authorization_code", code)
	data.Set("grant_type", "AUTHORIZATION_CODE")
	data.Set("redirect_url", *service.redirectUrl)
	data.Set("scope", *service.redirectUrl)

	encoded := data.Encode()
	body := strings.NewReader(encoded)

	req, err := http.NewRequest(http.MethodPost, tokenUrl, body)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Content-Length", strconv.Itoa(len(encoded)))
	req.Header.Set("Accept", "application/json")
	if err != nil {
		return nil, errortools.ErrorMessage(err)
	}

	return req, nil
}

func NewService(cfg *ServiceConfig) (*Service, *errortools.Error) {
	if cfg == nil {
		return nil, errortools.ErrorMessage("ServiceConfig must not be a nil pointer")
	}

	if cfg.ClientId == "" {
		return nil, errortools.ErrorMessage("ClientId not provided")
	}

	redirectUrl := defaultRedirectUrl
	if cfg.RedirectUrl != nil {
		redirectUrl = *cfg.RedirectUrl
	}

	var service = Service{
		clientId:     cfg.ClientId,
		clientSecret: cfg.ClientSecret,
		redirectUrl:  cfg.RedirectUrl,
	}

	//var getTokenRequestFunc = service.getTokenRequest

	oauth2ServiceConfig := oauth2.ServiceConfig{
		ClientId:        cfg.ClientId,
		ClientSecret:    cfg.ClientSecret,
		RedirectUrl:     redirectUrl,
		AuthUrl:         authUrl,
		TokenUrl:        tokenUrl,
		TokenHttpMethod: tokenHttpMethod,
		TokenSource:     cfg.TokenSource,
		//GetTokenFromRequestFunc: &getTokenRequestFunc,
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

	for {
		request, response, e := service.oAuth2Service.HttpRequest(requestConfig)
		if response != nil {
			if response.StatusCode == http.StatusTooManyRequests {
				reset := response.Header.Get("x-ratelimit-reset")
				resetInt, err := strconv.ParseInt(reset, 10, 64)
				if err == nil {
					if resetInt < 60*60 {
						time.Sleep(time.Duration(resetInt+1) + time.Second)
						continue
					}
				}
			}
		}
		if e != nil {
			if service.errorResponse.Message != "" {
				e.SetMessage(service.errorResponse.Message)
			}
		}

		if e != nil {
			return request, response, e
		}

		return request, response, nil
	}
}

func (service *Service) AuthorizeUrl(scope string, state string) string {
	if service.redirectUrl == nil {
		return ""
	}
	return fmt.Sprintf("%s?applicationId=%s&redirectUrl=%s&state=%s&scope=%s", authUrl, service.clientId, *service.redirectUrl, state, scope)
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
