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

type UsersNodesDataStruct struct {
	Data struct {
		Block struct {
			Nodes    []UserStruct   `json:"nodes"`
			PageInfo pageInfoStruct `json:"pageInfo"`
		} `json:"users"`
	} `json:"data"`
}

type UserStruct struct {
	Id                 string             `json:"id"`
	UserName           string             `json:"username,omitempty"`
	Name               string             `json:"name,omitempty"`
	Bot                bool               `json:"bot"`
	State              string             `json:"state,omitempty"`
	ProjectMemberships ProjectMemberships `json:"projectMemberships"`
	GroupMemberships   GroupMemberships   `json:"groupMemberships"`
}
type ProjectMemberships struct {
	Nodes []struct {
		AccessLevel AccessLevelStruct `json:"accessLevel"`
		Project     struct {
			NameWithNamespace string `json:"nameWithNamespace"`
			FullPath          string `json:"fullPath"`
		} `json:"project"`
	} `json:"nodes"`
}

type GroupMemberships struct {
	Nodes []struct {
		AccessLevel AccessLevelStruct `json:"accessLevel"`
		Group       struct {
			FullName string `json:"fullName"`
			FullPath string `json:"fullPath"`
		} `json:"group"`
	} `json:"nodes"`
}

type AccessLevelStruct struct {
	StringValue  string `json:"stringValue"`
	IntegerValue int    `json:"integerValue,omitempty"`
}

type UserInfoGraphQL struct {
	Data struct {
		Users struct {
			Nodes []struct {
				Id                 string `json:"id"`
				Name               string `json:"name"`
				Username           string `json:"username"`
				Bot                bool   `json:"bot"`
				State              string `json:"state"`
				ProjectMemberships struct {
				} `json:"projectMemberships"`
				GroupMemberships struct {
					Nodes []struct {
						AccessLevel struct {
							StringValue string `json:"stringValue"`
						} `json:"accessLevel"`
						Group struct {
							FullName string `json:"fullName"`
						} `json:"group"`
					} `json:"nodes"`
				} `json:"groupMemberships"`
			} `json:"nodes"`
			PageInfo struct {
				EndCursor   string `json:"endCursor"`
				HasNextPage bool   `json:"hasNextPage"`
			} `json:"pageInfo"`
		} `json:"users"`
	} `json:"data"`
}

func GetUsersInfo(c *gin.Context, cfg conf.Config, searchPattern string) ([]UserStruct, string) {

	tokenBlock, CSRFPageCookies, err := parsers.GetCSRFTokenPage(c, cfg)
	internal.CheckError(err)
	var GraphQLHeadersFromPage map[string]string
	err = json.Unmarshal([]byte(tokenBlock), &GraphQLHeadersFromPage)
	internal.CheckError(err)
	var ContentType string
	var PageId string
	var UsersInfoData UsersNodesDataStruct
	var UserNodes []UserStruct
	MaxElements := 20
	//log.Printf("SearchPattern: %s\n", searchPattern)
	if len(searchPattern) > 0 {
		MaxElements = 100
	} else {
		uinfo, _, err := api.GetUserInfo(c, cfg)
		if err != nil {
			log.Printf("%v", err)
		}
		searchPattern = uinfo.Username
	}
	CookiePack := api.GetCookiePack(c, CSRFPageCookies)
	internal.CheckError(err)
	for {
		jq := fmt.Sprintf(`
		query { users (after: "%s", search: "%s", first: %d, sort: CREATED_ASC) {
					nodes { 
						id name username bot state
						projectMemberships { nodes { accessLevel { stringValue } project { nameWithNamespace fullPath } } }
      					groupMemberships { nodes { accessLevel { stringValue } group { fullName fullPath} } } 
					}
					pageInfo { endCursor hasNextPage }
			  }
		}
		`, PageId, strings.ReplaceAll(searchPattern, "*", ""), MaxElements)

		jq = strings.ReplaceAll(jq, "\t", "")
		graphqlresp, bodyBytes, err := GraphQLData(c, cfg, GraphQLHeadersFromPage, CookiePack, jq, []string{})
		//fmt.Printf("Bytes: %v", string(bodyBytes))
		if err != nil {
			break
		}
		ContentType = graphqlresp.Header.Get("Content-Type")
		hasNextPage, NewPageId := CheckNextPage("users", bodyBytes)
		log.Printf("hasNextPage: %s NewPageId: %s\n", hasNextPage, NewPageId)

		err = json.Unmarshal(bodyBytes, &UsersInfoData)
		for _, node := range UsersInfoData.Data.Block.Nodes {
			UserNodes = append(UserNodes, node)
		}
		if !hasNextPage {
			break
		}
		if searchPattern == "" {
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
	return UserNodes, ContentType
}
