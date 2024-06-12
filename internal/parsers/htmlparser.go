package parsers

import (
	"errors"
	"github.com/aborche/gitlab-menu-extender-demo/internal/api"
	"github.com/aborche/gitlab-menu-extender-demo/internal/conf"
	"github.com/gin-gonic/gin"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
	"io"
	"log"
	"net/http"
	"slices"
	"strings"
)

func getAttribute(n *html.Node, key string) (string, error) {
	for _, attr := range n.Attr {
		if attr.Key == key {
			return attr.Val, nil
		}
	}
	return "", errors.New(key + " not exist in attribute!")
}
func copyAndSkipAttribute(n *html.Node, attrs []string) ([]html.Attribute, error) {
	var Attr []html.Attribute
	for _, attr := range n.Attr {
		if !slices.Contains(attrs, attr.Key) {
			Attr = append(Attr, attr)
		}
	}
	return Attr, nil
}
func GetSideBarMenu(n *html.Node) (string, error) {
	var sidebar string
	if n.Data == "aside" {
		sidebar, _ = getAttribute(n, "data-sidebar")
		if len(sidebar) > 0 {
			return sidebar, nil
		}
	}
	for child := n.FirstChild; child != nil; child = child.NextSibling {
		sidebar, _ = GetSideBarMenu(child)
		if len(sidebar) > 0 {
			return sidebar, nil
		}
	}
	return sidebar, nil
}

func ExtractCSRFToken(n *html.Node) (string, error) {
	var graphqlsettings string
	var err error
	if n.Data == "div" {
		divid, _ := getAttribute(n, "id")
		if divid == "graphiql-container" {
			graphqlsettings, err = getAttribute(n, "data-headers")
			if err != nil {
				return "", err
			}
			if len(graphqlsettings) > 0 {
				return graphqlsettings, nil
			}
		}
	}
	for child := n.FirstChild; child != nil; child = child.NextSibling {
		graphqlsettings, err = ExtractCSRFToken(child)
		if err != nil {
			return "", err
		}
		if len(graphqlsettings) > 0 {
			return graphqlsettings, nil
		}
	}
	return graphqlsettings, nil
}

func GetCSRFTokenPage(c *gin.Context, cfg conf.Config) (string, []*http.Cookie, error) {
	GraphQLPageStartPage, err := api.GetCustomResponse(c, cfg, "GET", "/-/graphql-explorer", nil)
	if err != nil {
		return "", nil, err
	}
	CSRFPageCookies := GraphQLPageStartPage.Cookies()
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Fatalln(err.Error())
		}
	}(GraphQLPageStartPage.Body)
	doc, err := html.Parse(GraphQLPageStartPage.Body)
	if err != nil {
		return "", nil, err
	}
	tokenBlock, err := ExtractCSRFToken(doc)
	if err != nil {
		return "", nil, err
	}
	return tokenBlock, CSRFPageCookies, err
}

func ReparseJSSettings(n *html.Node) (map[string]string, error) {
	settingsmap := make(map[string]string)
	jssettings := GetPageJSSettings([]string{}, n)
	if len(jssettings) < 1 {
		return nil, errors.New("JS Settings not found. Check gitlab URL.")
		//c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"status": "error", "message":})
		//return
	}
	for _, setpart := range jssettings {
		for _, part := range strings.Split(strings.Trim(setpart, "\n\r"), ";gon.") {
			if strings.Contains(part, "=") && !strings.Contains(part, "<!") && !strings.HasPrefix(part, "//") {
				set := strings.Split(part, "=")
				settingsmap[set[0]] = strings.Trim(set[1], `"`)
			}
		}
	}
	return settingsmap, nil
}

func GetPageJSSettings(settings []string, n *html.Node) []string {
	if n.Data == "script" {
		if n.FirstChild != nil {
			if n.FirstChild.Type == html.TextNode {
				if strings.Contains(n.FirstChild.Data, "window.gon={};") {
					settings = append(settings, n.FirstChild.Data)
				}
			}
		}
	}
	for child := n.FirstChild; child != nil; child = child.NextSibling {
		settings = GetPageJSSettings(settings, child)
	}
	return settings
}

func Reformat(n *html.Node) {
	if n.Type == html.ElementNode {
		if n.Data == "main" {
			newNode := &html.Node{
				Type:     html.ElementNode,
				Data:     "main",
				DataAtom: atom.Body,
				Attr: []html.Attribute{
					{Key: "class", Val: "content"},
					{Key: "id", Val: "content-body"},
				},
			}
			n.Parent.AppendChild(newNode)
			n.Parent.LastChild.AppendChild(&html.Node{
				Type: html.TextNode,
				Data: `{{ .BodyPage }}`, //| unescapeHTML
			})
			n.Parent.RemoveChild(n)
		} else if n.Data == "title" {
			if n.FirstChild.Type == html.TextNode {
				n.FirstChild.Data = "{{ .Title }}"
			}
		} else if n.Data == "nav" {
			n.Parent.AppendChild(&html.Node{
				Type:     html.ElementNode,
				Data:     "navbar",
				DataAtom: atom.Body,
				Attr: []html.Attribute{
					{Key: "id", Val: "navbar"},
				},
			})
			n.Parent.LastChild.AppendChild(&html.Node{
				Type: html.TextNode,
				//Data: `{{ unescapeHTML .Navigation  }}`,
				Data: `{{ .Navigation  }}`,
			})
			n.Parent.RemoveChild(n)
		} else if n.Data == "aside" {
			n.Attr, _ = copyAndSkipAttribute(n, []string{"data-sidebar"})
			n.Attr = append(n.Attr, html.Attribute{Key: "data-sidebar", Val: `{{ $.Sidebar }}`})
		}

	}
	for child := n.FirstChild; child != nil; child = child.NextSibling {
		Reformat(child)
	}
}
