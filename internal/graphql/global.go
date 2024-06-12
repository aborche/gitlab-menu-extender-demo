package graphql

import (
	"bytes"
	"encoding/json"
	"github.com/aborche/gitlab-menu-extender-demo/internal"
	"github.com/aborche/gitlab-menu-extender-demo/internal/api"
	"github.com/aborche/gitlab-menu-extender-demo/internal/conf"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
)

type GraphQL struct {
	Query     string   `json:"query"`
	Variables []string `json:"variables,omitempty"`
}

type pageInfoStruct struct {
	StartCursor     string `json:"startCursor"`
	EndCursor       string `json:"endCursor"`
	HasNextPage     bool   `json:"hasNextPage"`
	HasPreviousPage bool   `json:"hasPreviousPage"`
}

func MakeCustomRequest(c *gin.Context, cfg conf.Config, method string, path string, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequest(method, api.GetCtxBaseURL(c, cfg)+path, body)
	if err != nil {
		return nil, err
	}

	return req, err
}

func CheckNextPage(t string, bodyBytes []byte) (hasNextPage bool, nextPageId string) {
	switch t {
	case "users":
		var UsersPage UsersNodesDataStruct
		err := json.Unmarshal(bodyBytes, &UsersPage)
		internal.CheckError(err)
		return UsersPage.Data.Block.PageInfo.HasNextPage, UsersPage.Data.Block.PageInfo.EndCursor
	case "groups":
		var GroupsPage GroupsStruct
		err := json.Unmarshal(bodyBytes, &GroupsPage)
		internal.CheckError(err)
		return GroupsPage.Data.Groups.PageInfo.HasNextPage, GroupsPage.Data.Groups.PageInfo.EndCursor
	case "projects":
		var ProjectsPage ProjectsNodesDataStruct
		err := json.Unmarshal(bodyBytes, &ProjectsPage)
		internal.CheckError(err)
		return ProjectsPage.Data.Block.PageInfo.HasNextPage, ProjectsPage.Data.Block.PageInfo.EndCursor
	default:
		return false, ""
	}
}
func GraphQLData(c *gin.Context, cfg conf.Config, headers map[string]string, Cookies []*http.Cookie, query string, variables []string) (*http.Request, []byte, error) {
	//csrfcontent, err := helpers.GetCustomResponse(c, cfg, "GET", "/-/graphql-explorer", nil)
	opt := &GraphQL{
		Query:     query,
		Variables: []string{},
	}
	body, err := json.Marshal(opt)
	if err != nil {
		return nil, nil, err
	}

	graphqlreq, err := MakeCustomRequest(c, cfg, "POST", "/api/graphql", bytes.NewReader(body))
	if err != nil {
		return nil, nil, err
	}
	for _, cookie := range Cookies {
		graphqlreq.AddCookie(cookie)
	}

	for key, value := range headers {
		graphqlreq.Header.Set(key, value)
	}

	client := api.GetRedirectClient()
	graphqlresp, err := client.Do(graphqlreq)
	if err != nil {
		return nil, nil, err
	}
	defer graphqlresp.Body.Close()
	bodyBytes, err := io.ReadAll(graphqlresp.Body)
	if err != nil {
		return nil, nil, err
	}
	return graphqlreq, bodyBytes, err
}
