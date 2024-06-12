package gitlabmenu

import (
	"encoding/json"
	"github.com/aborche/gitlab-menu-extender-demo/internal"
	"github.com/aborche/gitlab-menu-extender-demo/internal/conf"
	"github.com/gin-gonic/gin"
	"io"
	"os"
	"strings"
)

func BuildTabMenu(c *gin.Context, cfg conf.Config) []TabMenu {
	tabs := make([]TabMenu, 0)
	splitStages := strings.Split(cfg.Stages, ",")
	for _, value := range splitStages {
		stage := strings.Split(value, ":")[0]
		stagingArea := strings.Split(value, ":")[1]
		tabPath := cfg.UrlPrefix + `deploy/` + stagingArea
		var tabStyle string
		if tabPath == c.Request.URL.Path {
			tabStyle = "active gl-tab-nav-item-active btn btn-block" //
		} else {
			tabStyle = ""
		}
		tabs = append(tabs, TabMenu{
			Name:        stage,
			Branch:      stagingArea,
			Path:        tabPath,
			ActiveStyle: tabStyle,
		})
	}
	return tabs
}

func BuildSidebarMenuFromFile(c *gin.Context, cfg conf.Config) []CurrentMenuItem {
	jsonFile, err := os.Open(cfg.MenuFile)
	internal.CheckError(err)
	defer jsonFile.Close()
	byteValue, _ := io.ReadAll(jsonFile)
	var CMenu []CurrentMenuItem
	err = json.Unmarshal(byteValue, &CMenu)
	internal.CheckError(err)

	for e, data := range CMenu {
		var ActiveMenu bool
		if len(data.Items) > 0 {
			for i, item := range data.Items {
				if item.Link == c.Request.URL.Path {
					CMenu[e].Items[i].IsActive = true
					ActiveMenu = true
				}
			}
		}
		if ActiveMenu {
			CMenu[e].IsActive = true
		} else {
			if data.Link == c.Request.URL.Path {
				CMenu[e].IsActive = true
			}
		}
	}
	return CMenu
}

func BuildSidebarMenuFromStruct(c *gin.Context, cfg conf.Config) []CurrentMenuItem {
	var pillcount int
	pillcount = 0
	CMenu := []CurrentMenuItem{
		{
			ID:    "back",
			Title: "Back to Gitlab",
			Icon:  "go-back",
			Link:  "/",
		},
		{
			ID:    "helper",
			Title: "Deploy Information",
			Icon:  "information-o",
			Link:  cfg.UrlPrefix,
		},
		{
			ID:          "deploy",
			Title:       "Deploy Area",
			Icon:        "rocket-launch",
			Link:        cfg.UrlPrefix + "deploy",
			LinkClasses: "show",
			Items: []CurrentMenuItemItem{
				{
					ID:        "dev",
					Title:     "Deploy to DEV",
					Icon:      "dash-circle",
					Link:      cfg.UrlPrefix + "deploy/to-development",
					PillCount: &pillcount,
				},
				{
					ID:        "test",
					Title:     "Deploy to TEST",
					Icon:      "check-circle-dashed",
					Link:      cfg.UrlPrefix + "deploy/to-testing",
					PillCount: nil,
				},
				{
					ID:        "acc",
					Title:     "Deploy to ACC",
					Icon:      "check-circle",
					Link:      cfg.UrlPrefix + "deploy/to-acceptance",
					PillCount: nil,
				},
				{
					ID:        "prod",
					Title:     "Deploy to PROD",
					Icon:      "check-circle-filled",
					Link:      cfg.UrlPrefix + "deploy/to-production",
					PillCount: nil,
				},
			},
		},
	}

	for e, data := range CMenu {
		var ActiveMenu bool
		if len(data.Items) > 0 {
			for i, item := range data.Items {
				if item.Link == c.Request.URL.Path {
					CMenu[e].Items[i].IsActive = true
					ActiveMenu = true
				}
			}
		}
		if ActiveMenu {
			CMenu[e].IsActive = true
		} else {
			if data.Link == c.Request.URL.Path {
				CMenu[e].IsActive = true
			}
		}
	}
	return CMenu
}
