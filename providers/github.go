package providers

import (
	"context"
	"net/http"
	"strings"

	"github.com/drone/drone-go/drone"
	"github.com/drone/go-scm/scm"
	"github.com/drone/go-scm/scm/driver/github"
	"github.com/drone/go-scm/scm/transport"

	"github.com/prometheus/client_golang/prometheus"
)

func GetGithubFilesChanged(repo drone.Repo, build drone.Build, token string) ([]string, error) {
	newctx := context.Background()
	client := github.NewDefault()
	client.Client = &http.Client{
		Transport: &transport.BearerToken{
			Token: token,
		},
	}

	var changes []*scm.Change
	var result *scm.Response
	var err error

	if build.Before == "" || build.Before == scm.EmptyCommit {
		changes, result, err = client.Git.ListChanges(newctx, strings.Join([]string{repo.Namespace, repo.Name}, "/"), build.After, scm.ListOptions{})
		if err != nil {
			return nil, err
		}
	} else {
		changes, result, err = client.Git.CompareChanges(newctx, strings.Join([]string{repo.Namespace, repo.Name}, "/"), build.Before, build.After, scm.ListOptions{})
		if err != nil {
			return nil, err
		}
	}

	GithubApiCount.Set(float64(result.Rate.Remaining))

	var files []string
	for _, c := range changes {
		files = append(files, c.Path)
	}

	return files, nil
}

var (
	GithubApiCount = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "github_api_calls_remaining",
			Help: "Total number of github api calls per hour remaining",
		})
)

func init() {
	prometheus.MustRegister(GithubApiCount)
}
