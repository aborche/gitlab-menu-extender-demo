package api

import (
	"fmt"
	"github.com/aborche/gitlab-menu-extender-demo/internal/conf"
	"github.com/aborche/go-gitlab"
	"github.com/gin-gonic/gin"
	"log"
)

func CheckAuth(c *gin.Context, cfg conf.Config) (StatusCode int, err error) {
	_, resp, err := GetUserInfo(c, cfg)
	return resp.StatusCode, err
}

func GetUserInfo(c *gin.Context, cfg conf.Config) (*gitlab.User, *gitlab.Response, error) {
	git, err := Client(c, cfg)
	if err != nil {
		return nil, nil, err
	}
	users, resp, err := git.Users.CurrentUser()
	if err != nil {
		return nil, resp, err
	}
	return users, resp, err
}

func GetUserMemberShip(c *gin.Context, cfg conf.Config, id int) (any, error) {
	git, err := Client(c, cfg)
	if err != nil {
		return nil, err
	}
	opt := &gitlab.GetUserMembershipOptions{}
	membership, _, err := git.Users.GetUserMemberships(id, opt)
	log.Printf("%+v", membership)
	if err != nil {
		return nil, err
	}
	return membership, err
}

func GetUserIP(c *gin.Context) string {
	RemoteIP := c.RemoteIP()
	return RemoteIP
}

func GetHeaders(c *gin.Context) {
	fmt.Println("-=HEADERS=-")
	for name, headers := range c.Request.Header {
		for _, h := range headers {
			fmt.Printf("%v: %v\n", name, h)
		}
	}
}
