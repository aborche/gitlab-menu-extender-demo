package api

import (
	"log"
	"regexp"
	"strings"
)

func CleanURL(url string) string {
	var re = regexp.MustCompile(`^http(s)?://(.*)(/-/pipelines/|:)(.*)`)
	s := ""
	if re.MatchString(url) {
		s = re.ReplaceAllString(url, `$2/$4`)
	} else {
		s = strings.Replace(url, ":", "/", -1)
	}
	log.Printf("CleanUp: %v", s)
	return s
}
