package api

import (
	"errors"
	"fmt"
	"github.com/aborche/gitlab-menu-extender-demo/internal/conf"
	"github.com/aborche/go-gitlab"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
)

func Client(c *gin.Context, cfg conf.Config) (*gitlab.Client, error) {
	ctxBaseUrl := GetCtxBaseAPIURL(c, cfg)
	ctxSessionCookies := GetCtxCookies(c)
	git, err := gitlab.NewCookieClient(ctxSessionCookies, gitlab.WithBaseURL(ctxBaseUrl))
	if err != nil {
		return nil, fmt.Errorf("Failed to create client: %v", err)
	}
	if len(ctxSessionCookies) < 2 {
		return nil, errors.New("no requested cookies found")
	}
	return git, nil
}

func GetRedirectClient() *http.Client {
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	return client
}

func GetCustomResponse(c *gin.Context, cfg conf.Config, method string, path string, body io.Reader) (*http.Response, error) {
	client := GetRedirectClient()
	req, err := http.NewRequest(method, GetCtxBaseURL(c, cfg)+path, body)
	if err != nil {
		return nil, err
	}
	for _, cookie := range GetCtxCookies(c) {
		req.AddCookie(cookie)
	}

	resp, err := client.Do(req)
	return resp, err
}
