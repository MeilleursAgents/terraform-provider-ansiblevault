package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/go-github/v25/github"
	"github.com/pkg/errors"
	"golang.org/x/oauth2"
)

type app struct {
	owner  string
	repo   string
	client *github.Client
}

var (
	contextTimeout = time.Second * 10
	buildDirectory = "dist"
	owner          = "MeilleursAgents"
	repository     = "terraform-provider-ansiblevault"
)

func getClient(accessToken string) *github.Client {
	ctx, cancel := context.WithTimeout(context.Background(), contextTimeout)
	defer cancel()

	ts := oauth2.StaticTokenSource(
		&oauth2.Token{
			AccessToken: accessToken,
		},
	)

	client := github.NewClient(oauth2.NewClient(ctx, ts))

	return client
}

func new(accessToken string) *app {
	return &app{
		owner:  owner,
		repo:   repository,
		client: getClient(accessToken),
	}
}

func (a *app) publishArtifactsOnTag(tag string, buildAssets []string) error {

	log.Println("==> Retrieve release informations from tag ...")
	release, err := a.retrieveTag(tag)

	if err != nil {
		return errors.Wrapf(err, "When retrieve release information %s", release.GetHTMLURL())
	}

	log.Println("==> Retrieve release assets ...")
	releaseRemoteAssets, err := a.listAssets(release.GetID())

	if err != nil {
		return errors.Wrapf(err, "When retrieve assets %s", release.GetHTMLURL())
	}

	assetsToUpload := a.computeUploadAssets(releaseRemoteAssets, buildAssets)

	if len(assetsToUpload) > 0 {
		log.Println("==> Upload release assets ...")
		err = a.uploadAssets(release.GetID(), assetsToUpload)

		if err != nil {
			return errors.Wrapf(err, "uploadAssets on %s with %+q", release.GetHTMLURL(), assetsToUpload)
		}
	} else {
		log.Println("==> Nothing to upload ...")
	}

	return nil
}

func (a *app) computeUploadAssets(releaseRemoteAssets []*github.ReleaseAsset, buildAssets []string) []string {
	remoteAssetsNotFound := []string{}
	remoteAssets := mergeReleaseAssetsNames(releaseRemoteAssets)

	for _, value := range buildAssets {
		if !contains(remoteAssets, strings.Replace(value, fmt.Sprintf("%s/", buildDirectory), "", -1)) {
			remoteAssetsNotFound = append(remoteAssetsNotFound, value)
		}
	}

	return remoteAssetsNotFound
}

func (a *app) retrieveTag(tag string) (*github.RepositoryRelease, error) {
	ctx, cancel := context.WithTimeout(context.Background(), contextTimeout)
	defer cancel()
	release, _, err := a.client.Repositories.GetReleaseByTag(ctx, a.owner, a.repo, tag)

	if err != nil {
		return nil, err
	}

	return release, nil
}

func (a *app) listAssets(releaseID int64) ([]*github.ReleaseAsset, error) {
	ctx, cancel := context.WithTimeout(context.Background(), contextTimeout)
	defer cancel()

	result := []*github.ReleaseAsset{}
	page := 1

	for {
		assets, res, err := a.client.Repositories.ListReleaseAssets(ctx, a.owner, a.repo, releaseID, &github.ListOptions{Page: page})
		if err != nil {
			return nil, err
		}

		result = append(result, assets...)

		if res.NextPage <= page {
			break
		}

		page = res.NextPage
	}

	return result, nil
}

func (a *app) uploadAssets(releaseID int64, assets []string) error {
	ctx, cancel := context.WithTimeout(context.Background(), contextTimeout)
	defer cancel()

	for _, asset := range assets {
		log.Println(fmt.Sprintf("====> Upload %s", asset))

		opts := &github.UploadOptions{
			Name: filepath.Base(asset),
		}

		file, err := os.Open(asset)
		if err != nil {
			return errors.Wrapf(err, "failed to open file")
		}

		if _, _, err = a.client.Repositories.UploadReleaseAsset(ctx, a.owner, a.repo, releaseID, opts, file); err != nil {
			return errors.Wrapf(err, "failed to upload asset: %s", asset)
		}
	}

	return nil
}

func mergeReleaseAssetsNames(releaseRemoteAssets []*github.ReleaseAsset) []string {
	remoteAssets := []string{}

	for _, value := range releaseRemoteAssets {
		remoteAssets = append(remoteAssets, *value.Name)
	}

	return remoteAssets
}
