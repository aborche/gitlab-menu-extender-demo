package graphql

import (
	"encoding/json"
	"fmt"
	"github.com/aborche/gitlab-menu-extender-demo/internal"
	"github.com/aborche/gitlab-menu-extender-demo/internal/api"
	"github.com/aborche/gitlab-menu-extender-demo/internal/conf"
	"github.com/aborche/gitlab-menu-extender-demo/internal/parsers"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

type GroupsStruct struct {
	Data struct {
		Groups struct {
			Nodes    []any          `json:"nodes"`
			PageInfo pageInfoStruct `json:"pageInfo"`
		} `json:"groups"`
	} `json:"data"`
}

type GroupsDataStruct struct {
	Data struct {
		Groups struct {
			Nodes    []GroupsNodeStruct `json:"nodes"`
			PageInfo pageInfoStruct     `json:"pageInfo"`
		} `json:"groups"`
	} `json:"data"`
}

type GroupMemberNodeStruct struct {
	User        UserStruct        `json:"user"`
	AccessLevel AccessLevelStruct `json:"accessLevel"`
}

type GroupsNodeStruct struct {
	Id                string `json:"id"`
	Name              string `json:"name"`
	Description       string `json:"description"`
	FullName          string `json:"fullName"`
	FullPath          string `json:"fullPath"`
	GroupMembersCount int    `json:"groupMembersCount"`
	GroupMembers      struct {
		Nodes []GroupMemberNodeStruct `json:"nodes"`
	} `json:"groupMembers"`
}

func GetGroupsInfo(c *gin.Context, cfg conf.Config) ([]byte, string) {
	tokenBlock, CSRFPageCookies, err := parsers.GetCSRFTokenPage(c, cfg)
	internal.CheckError(err)
	var GraphQLHeadersFromPage map[string]string
	err = json.Unmarshal([]byte(tokenBlock), &GraphQLHeadersFromPage)
	internal.CheckError(err)
	var ContentType string
	var PageId string
	var GroupInfoData GroupsDataStruct
	var GroupNodes []GroupsNodeStruct
	//var GraphQLPage helpers.GraphQLStruct
	//var Edges []helpers.EdgesNode

	CookiePack := api.GetCookiePack(c, CSRFPageCookies)
	internal.CheckError(err)
	for {
		jq := fmt.Sprintf(`
			query {
			  groups (after: "%s") {
				nodes {
				  id
				  name
				  description
				  fullName
				  fullPath
				  groupMembersCount
				  groupMembers {
					nodes {
					  user {
						id
						username
						name
					  }
					  accessLevel {
						stringValue
						integerValue
					  }
					}
				  }
				}
				pageInfo {
				  endCursor
				  hasNextPage
				}
			  }
			}
			`, PageId)

		//fmt.Printf("JQ: %v", jq)
		//m := regexp.MustCompile("\t|\n|\r")
		//jq = m.ReplaceAllString(jq, " ")
		jq = strings.ReplaceAll(jq, "\t", "")
		graphqlresp, bodyBytes, err := GraphQLData(c, cfg, GraphQLHeadersFromPage, CookiePack, jq, []string{})
		//fmt.Printf("Bytes: %v", string(bodyBytes))
		if err != nil {
			break
		}
		ContentType = graphqlresp.Header.Get("Content-Type")
		hasNextPage, NewPageId := CheckNextPage("groups", bodyBytes)

		err = json.Unmarshal(bodyBytes, &GroupInfoData)
		for _, node := range GroupInfoData.Data.Groups.Nodes {
			GroupNodes = append(GroupNodes, node)
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
	data, err := json.Marshal(GroupNodes)
	internal.CheckError(err)

	return data, ContentType
}
