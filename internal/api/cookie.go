package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetCookiePack(c *gin.Context, CSRFCookies []*http.Cookie) []*http.Cookie {
	var CookiePack []*http.Cookie
	var GitlabSessionCookieFound bool
	for _, cookie := range CSRFCookies {
		fmt.Printf("CSRFCOOKIE: %v\n", cookie)
		CookiePack = append(CookiePack, cookie)
		if cookie.Name == "_gitlab_session" {
			GitlabSessionCookieFound = true
		}
	}
	for _, cookie := range GetCtxCookies(c) {
		//fmt.Printf("CTXCOOKIE: %v\n", cookie)
		if cookie.Name != "_gitlab_session" {
			CookiePack = append(CookiePack, cookie)
		} else if !GitlabSessionCookieFound {
			CookiePack = append(CookiePack, cookie)
		}
	}
	return CookiePack
}
