package providers

import (
	"context"
	"net/http"

	"github.com/drone/drone-go/drone"
	"github.com/drone/go-scm/scm"
	"github.com/drone/go-scm/scm/driver/gitea"
	"github.com/drone/go-scm/scm/transport"

	"github.com/prometheus/client_golang/prometheus"
)

type giteaDiffs struct {
	Diffs []struct {
		Destination struct {
			Name string `json:"toString"`
		} `json:"destination"`
	} `json:"diffs"`
}

func GetGiteaFilesChanged(repo drone.Repo, build drone.Build, token string, uri string) ([]string, error) {
	var err error

	newctx := context.Background()
	client, err := gitea.New(uri)
	if err != nil {
		return nil, err
	}
	client.Client = &http.Client{
		Transport: &transport.BearerToken{
			Token: token,
		},
	}

	var changes []*scm.Change
	var result *scm.Response

	if build.Before == "" || build.Before == scm.EmptyCommit {
		changes, result, err = client.Git.ListChanges(newctx, repo.Slug, build.After, scm.ListOptions{})
		if err != nil {
			return nil, err
		}
	} else {
		changes, result, err = client.Git.CompareChanges(newctx, repo.Slug, build.Before, build.After, scm.ListOptions{})
		if err != nil {
			return nil, err
		}
	}

	GiteaApiCount.Set(float64(result.Rate.Remaining))

	var files []string
	for _, c := range changes {
		files = append(files, c.Path)
	}

	return files, nil
}

var (
	GiteaApiCount = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "gitea_api_calls_remaining",
			Help: "Total number of gitea api calls per hour remaining",
		})
)

func init() {
	prometheus.MustRegister(GiteaApiCount)
}