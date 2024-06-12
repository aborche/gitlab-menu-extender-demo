package conf

import (
	"github.com/spf13/viper"
)

type Config struct {
	GitlabToken       string
	CustomGitlabURL   string
	GitlabUrlFromHost bool
	UrlPrefix         string
	MenuFile          string
	DefaultSection    string
	Stages            string
}

var (
	defaults = map[string]interface{}{
		"GitlabUrlFromHost": true,
		"UrlPrefix":         "/-/helper",
		"MenuFile":          "templates/sidebarmenushort.json",
		"DefaultSection":    "/dashboard/projects",
		"Stages":            "DEV:to-development,TEST:to-testing",
	}
)

func LoadConfig() (config Config, err error) {
	for k, v := range defaults {
		viper.SetDefault(k, v)
	}
	viper.SetEnvPrefix("helper")
	viper.BindEnv("GitlabToken")
	viper.BindEnv("CustomGitlabURL")
	viper.BindEnv("GitlabUrlFromHost")
	viper.BindEnv("UrlPrefix")
	viper.BindEnv("MenuFile")
	viper.BindEnv("DefaultSection")
	viper.BindEnv("Stages")

	err = viper.Unmarshal(&config)
	return

}
