package api

import (
	"github.com/aborche/gitlab-menu-extender-demo/internal/conf"
	"github.com/gin-gonic/gin"
	"net/http"
	"slices"
	"strings"
)

func GetCtxBaseAPIURL(c *gin.Context, config conf.Config) string {
	scheme := "http"
	if c.Request.Header.Get("X-Forwarded-Proto") != "" {
		scheme = c.Request.Header.Get("X-Forwarded-Proto")
	}
	GitlabAPIHost := c.Request.Host
	if !config.GitlabUrlFromHost {
		GitlabAPIHost = config.CustomGitlabURL
	}
	ctxBaseURL := scheme + "://" + GitlabAPIHost + "/api/v4/"
	return ctxBaseURL
}

func GetCtxBaseURL(c *gin.Context, config conf.Config) string {
	scheme := "http"
	if c.Request.Header.Get("X-Forwarded-Proto") != "" {
		scheme = c.Request.Header.Get("X-Forwarded-Proto")
	}
	GitlabAPIHost := c.Request.Host
	if !config.GitlabUrlFromHost {
		GitlabAPIHost = config.CustomGitlabURL
	}
	ctxBaseURL := scheme + "://" + GitlabAPIHost
	return ctxBaseURL
}

func GetCtxCookieCheckURL(c *gin.Context, config conf.Config) string {
	scheme := "http"
	if c.Request.Header.Get("X-Forwarded-Proto") != "" {
		scheme = c.Request.Header.Get("X-Forwarded-Proto")
	}
	GitlabAPIHost := c.Request.Host
	if !config.GitlabUrlFromHost {
		GitlabAPIHost = config.CustomGitlabURL
	}
	ctxBaseURL := scheme + "://" + GitlabAPIHost + "/" + strings.TrimLeft(config.DefaultSection, "/")
	return ctxBaseURL
}

func GetCtxCookies(c *gin.Context) (GitlabCookies []*http.Cookie) {
	GitlabAuthCookies := []string{"known_sign_in", "_gitlab_session", "remember_user_token"}
	for _, value := range c.Request.Cookies() {
		if slices.Contains(GitlabAuthCookies, value.Name) {
			GitlabCookies = append(GitlabCookies, value)
		}
	}
	return GitlabCookies
}
