package main

import (
	"flag"
	"log"
	"strings"
	"time"
)

var (
	targets = map[string][]string{
		"linux": []string{
			"amd64",
			"386",
		},
		"darwin": []string{
			"amd64",
		},
	}
)

const (
	gitRepoURL = "github.com/MeilleursAgents/terraform-provider-ansiblevault"
)

func main() {
	flagAccessToken := flag.String("access-token", "", "Github Access Token for uploading assets")
	flag.Parse()

	log.Println("[GOLANG] Build")

	tag, err := executeCmd([]string{}, "git", "tag", "--sort=-taggerdate")

	if err != nil {
		log.Fatal(err)
	}

	date := time.Now().Format("2006-02-01_03:04:05PM")
	buildAssets, err := build(date, tag)

	if err != nil {
		log.Fatalf("%+v", err)
	}

	log.Println("[GITHUB] Upload")

	app := new(*flagAccessToken)
	err = app.publishArtifactsOnTag(tag, buildAssets)

	if err != nil {
		log.Fatalf("%+v", err)
	}
}

func contains(a []string, s string) bool {
	for _, n := range a {
		if strings.EqualFold(s, n) {
			return true
		}
	}
	return false
}
