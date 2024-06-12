package graphql

import (
	"encoding/json"
	"fmt"
	"github.com/aborche/gitlab-menu-extender-demo/internal"
	"github.com/aborche/gitlab-menu-extender-demo/internal/api"
	"github.com/aborche/gitlab-menu-extender-demo/internal/conf"
	"github.com/aborche/gitlab-menu-extender-demo/internal/parsers"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strings"
)

type ProjectsNodesDataStruct struct {
	Data struct {
		Block struct {
			Nodes    []any          `json:"nodes"`
			PageInfo pageInfoStruct `json:"pageInfo"`
		} `json:"projects"`
	} `json:"data"`
}

type ProjectsDataStruct struct {
	Data struct {
		Projects struct {
			Nodes    []ProjectsNodeStruct `json:"nodes"`
			PageInfo pageInfoStruct       `json:"pageInfo"`
		} `json:"projects"`
	} `json:"data"`
}
type ProjectsNodeStruct struct {
	ID             string `json:"id"`
	Name           string `json:"name"`
	FullPath       string `json:"fullPath"`
	ProjectMembers struct {
		Nodes []struct {
			User        UserStruct        `json:"user"`
			AccessLevel AccessLevelStruct `json:"accessLevel"`
		} `json:"nodes"`
	} `json:"projectMembers"`
}

func GetProjectsMembersInfo(c *gin.Context, cfg conf.Config, searchPattern string) ([]ProjectsNodeStruct, string) {
	tokenBlock, CSRFPageCookies, err := parsers.GetCSRFTokenPage(c, cfg)
	internal.CheckError(err)
	var GraphQLHeadersFromPage map[string]string
	err = json.Unmarshal([]byte(tokenBlock), &GraphQLHeadersFromPage)
	internal.CheckError(err)
	var ContentType string
	var PageId string

	var ProjectsNodes []ProjectsNodeStruct
	Count := 0
	MaxElements := 50
	if len(searchPattern) > 0 {
		MaxElements = 100
	}
	CookiePack := api.GetCookiePack(c, CSRFPageCookies)
	internal.CheckError(err)
	for {
		var ProjectsInfoData ProjectsDataStruct
		jq := fmt.Sprintf(`
			query {
			  projects(membership: true, search: "%s", first: %d, after: "%s", sort: "name_asc") {
				nodes { id name fullPath
				  projectMembers (sort: ACCESS_LEVEL_ASC) {
					  nodes {
						user { id name username bot state }
						accessLevel { stringValue }
					  }
				  }
				}
				pageInfo { endCursor hasNextPage }
			  }
			}
		`, strings.ReplaceAll(searchPattern, "*", ""), MaxElements, PageId)
		jq = strings.ReplaceAll(jq, "\t", "")
		graphqlresp, bodyBytes, err := GraphQLData(c, cfg, GraphQLHeadersFromPage, CookiePack, jq, []string{})
		if err != nil {
			break
		}
		ContentType = graphqlresp.Header.Get("Content-Type")
		hasNextPage, NewPageId := CheckNextPage("projects", bodyBytes)
		fmt.Printf("NextPage: %v => NextPageId: %v\n", hasNextPage, NewPageId)
		err = json.Unmarshal(bodyBytes, &ProjectsInfoData)
		log.Printf("ProjectsCount: %d", len(ProjectsInfoData.Data.Projects.Nodes))
		for _, node := range ProjectsInfoData.Data.Projects.Nodes {
			Count += 1
			ProjectsNodes = append(ProjectsNodes, node)
		}
		if searchPattern == "" {
			break
		}
		if !hasNextPage {
			break
		}

		if NewPageId != PageId {
			PageId = NewPageId
		}
	}

	for _, cookie := range CookiePack {
		c.SetSameSite(http.SameSiteStrictMode)
		c.SetCookie(cookie.Name, cookie.Value, cookie.MaxAge, cookie.Path, cookie.Domain, cookie.Secure, cookie.HttpOnly)
	}
	return ProjectsNodes, ContentType
}
