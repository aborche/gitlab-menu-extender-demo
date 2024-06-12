package renders

import (
	"bytes"
	"encoding/json"
	"github.com/aborche/gitlab-menu-extender-demo/internal"
	"github.com/aborche/gitlab-menu-extender-demo/internal/api"
	"github.com/aborche/gitlab-menu-extender-demo/internal/conf"
	"github.com/aborche/gitlab-menu-extender-demo/internal/gitlabmenu"
	"github.com/aborche/gitlab-menu-extender-demo/internal/graphql"
	"github.com/aborche/gitlab-menu-extender-demo/internal/parsers"
	"github.com/gin-gonic/gin"
	"golang.org/x/net/html"
	"html/template"
	"io"
	"log"
	"net/http"
)

func RenderProjects(c *gin.Context, cfg conf.Config, funcMap template.FuncMap, searchPattern string) (bytes.Buffer, error) {
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	req, err := http.NewRequest("GET", api.GetCtxCookieCheckURL(c, cfg), nil)

	internal.CheckError(err)
	for _, cookie := range api.GetCtxCookies(c) {
		req.AddCookie(cookie)
	}

	// Get gitlab start page
	resp, err := client.Do(req)
	internal.CheckError(err)
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		c.DataFromReader(http.StatusFound,
			-1,
			resp.Header.Get("content-type"),
			resp.Body,
			map[string]string{
				"Location": api.GetCtxBaseURL(c, cfg),
			})
		return bytes.Buffer{}, nil
	}
	// Parse start page into html.Node
	doc, err := html.Parse(resp.Body)
	internal.CheckError(err)

	// Extract gon.settings
	JSSettings, err := parsers.ReparseJSSettings(doc)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"status": "error", "message": err.Error()})
		return bytes.Buffer{}, err
	}

	// Transform sidebar
	NewSideBarJson, err := gitlabmenu.TransformSidebar(c, cfg, doc, true)
	internal.CheckError(err)

	// Define icons path from gon.settings
	funcMap["getIconsBundle"] = func() template.HTML { return template.HTML(JSSettings["sprite_icons"]) }

	// Reformat original document
	parsers.Reformat(doc)
	// Render doc into html.Template
	var b bytes.Buffer
	_ = io.Writer(&b)
	err = html.Render(&b, doc)
	internal.CheckError(err)
	// Recreate new template from changed start page
	t := template.Must(template.New("").Funcs(funcMap).Parse(b.String()))

	tmpl, err := template.
		New("helper_projects.tmpl").
		Funcs(funcMap).
		ParseFiles("templates/helper_projects.tmpl")

	internal.CheckError(err)

	// Get userinfo from gitlab api
	ApiUserInfo, _, err := api.GetUserInfo(c, cfg)
	if err != nil {
		log.Printf("%v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return bytes.Buffer{}, err
	}

	// Buffering existing data into rendered data block
	var buf bytes.Buffer
	_ = io.Writer(&buf)
	UserInfo, _ := json.Marshal(ApiUserInfo)
	ProjectsInfo, _ := graphql.GetProjectsMembersInfo(c, cfg, searchPattern)
	err = tmpl.Execute(&buf, map[string]interface{}{
		"tbl":      ProjectsInfo,
		"UserInfo": string(UserInfo),
		"gon":      JSSettings,
	})
	internal.CheckError(err)

	// Create map for final template render
	data := map[string]interface{}{
		"Title":      "Projects Audit",
		"Sidebar":    string(NewSideBarJson),
		"Navigation": `Projects Audit`,
		"BodyPage":   template.HTML(buf.String()),
	}

	// Render all
	var bf bytes.Buffer
	_ = io.Writer(&bf)
	err = t.Execute(&bf, data)
	internal.CheckError(err)
	return bf, nil
}
