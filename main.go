package main

import (
	"encoding/json"
	"github.com/aborche/gitlab-menu-extender-demo/internal"
	"github.com/aborche/gitlab-menu-extender-demo/internal/api"
	"github.com/aborche/gitlab-menu-extender-demo/internal/conf"
	"github.com/aborche/gitlab-menu-extender-demo/internal/graphql"
	"github.com/aborche/gitlab-menu-extender-demo/internal/renders"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func main() {
	cfg, err := conf.LoadConfig()
	if err != nil {
		log.Fatalf("Error reading config file: %s", err)
	}
	log.Printf("%+v", cfg)
	if !cfg.GitlabUrlFromHost && cfg.CustomGitlabURL == "" {
		log.Fatalf("Error in config file. When GitlabUrlFromHost is false, CustomGitlabURL must be defined")
	}

	r := gin.Default()
	funcMap := renders.SetFuncMaps(cfg)

	r.SetFuncMap(funcMap)
	rg := r.Group(cfg.UrlPrefix)
	r.LoadHTMLGlob("templates/*.tmpl")

	rg.GET("/getUserInfo", func(c *gin.Context) {
		StatusCode, AuthError := api.CheckAuth(c, cfg)
		if err != nil {
			c.JSON(StatusCode, gin.H{"error": AuthError.Error()})
		}
		//password := c.Param("password")
		uinfo, _, err := api.GetUserInfo(c, cfg)
		if err != nil {
			log.Printf("%v", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, any(uinfo))
	})

	rg.GET("/", func(c *gin.Context) {
		searchPattern := c.DefaultQuery("search", "")
		bf, err := renders.RenderMain(c, cfg, funcMap, searchPattern)
		internal.CheckError(err)
		c.Data(http.StatusOK, "text/html; charset=utf-8", bf.Bytes())
	})

	rg.GET("/graph", func(c *gin.Context) {
		UserInfo, ContentType := graphql.GetUsersInfo(c, cfg, "")
		data, err := json.Marshal(UserInfo)
		internal.CheckError(err)
		c.Data(http.StatusOK, ContentType, data)
		return
	})

	rg.GET("/projects", func(c *gin.Context) {
		searchPattern := c.DefaultQuery("search", "")
		bf, err := renders.RenderProjects(c, cfg, funcMap, searchPattern)
		internal.CheckError(err)
		c.Data(http.StatusOK, "text/html; charset=utf-8", bf.Bytes())
		return
	})

	rg.GET("/users", func(c *gin.Context) {
		searchPattern := c.DefaultQuery("search", "")
		bf, err := renders.RenderUsers(c, cfg, funcMap, searchPattern)
		internal.CheckError(err)
		c.Data(http.StatusOK, "text/html; charset=utf-8", bf.Bytes())
		return
	})

	r.Run()
}
