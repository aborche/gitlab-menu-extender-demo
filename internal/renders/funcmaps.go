package renders

import (
	"encoding/json"
	"github.com/aborche/gitlab-menu-extender-demo/internal/conf"
	"github.com/jackc/pgx/v5/pgtype"
	"html/template"
	"path/filepath"
	"time"
)

func SetFuncMaps(cfg conf.Config) template.FuncMap {
	funcMap := template.FuncMap{}
	funcMap["unescapeHTML"] = func(s string) template.HTML { return template.HTML(s) }

	funcMap["formatDate"] = func(s pgtype.Date) string {
		loc, _ := time.LoadLocation("Europe/Moscow")
		//return `<span class="gl-text-red-500">` + s.Time.In(loc).Format(time.RFC822) + `</span>`
		return s.Time.In(loc).Format(time.RFC822)
	}

	funcMap["marshal"] = func(v interface{}) template.JS {
		a, _ := json.Marshal(v)
		return template.JS(a)
	}

	funcMap["cutId"] = func(v string) string {
		Id := filepath.Base(v)
		return Id
	}

	funcMap["urlPrefix"] = func() string {
		return cfg.UrlPrefix
	}

	funcMap["getIconsBundle"] = func() template.HTML { return "" }

	return funcMap
}
