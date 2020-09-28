package plugin

import (
	"context"
	"fmt"
	"github.com/ktrysmt/go-bitbucket"

	"github.com/drone/drone-go/drone"
	"github.com/google/go-github/github"
	_ "github.com/ktrysmt/go-bitbucket"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
)

func getFilesChanged(repo drone.Repo, build drone.Build, token string, provider string) ([]string, error) {
	if provider == "github" {
		changedFiles, err := getGithubFilesChanged(repo, build, token)
		if err != nil {
			return nil, err
		}
		return changedFiles, nil
	} else if provider == "bitbucket" {
		changedFiles, err := getBBFilesChanged(repo, build, token)
		if err != nil {
			return nil, err
		}
		return changedFiles, nil
	}
	return nil, fmt.Errorf("unsupported provider %s", provider)
}

func getGithubFilesChanged(repo drone.Repo, build drone.Build, token string) ([]string, error) {
	newctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(newctx, ts)

	client := github.NewClient(tc)

	var commitFiles []github.CommitFile
	if build.Before == "" || build.Before == "0000000000000000000000000000000000000000" {
		response, _, err := client.Repositories.GetCommit(newctx, repo.Namespace, repo.Name, build.After)
		if err != nil {
			return nil, err
		}
		commitFiles = response.Files
	} else {
		response, _, err := client.Repositories.CompareCommits(newctx, repo.Namespace, repo.Name, build.Before, build.After)
		if err != nil {
			return nil, err
		}
		commitFiles = response.Files
	}
	rateLimit, _, err := client.RateLimits(newctx)
	if err != nil {
		logrus.Fatalln("No metrics")
	}
	//metrics.GithubApiCount.Set(float64(rateLimit.Core.Remaining))
	GithubApiCount.Set(float64(rateLimit.Core.Remaining))
	var files []string
	for _, f := range commitFiles {
		files = append(files, *f.Filename)
	}

	return files, nil
}

func getBBFilesChanged(repo drone.Repo, build drone.Build, token string) ([]string, error) {
	var files []string
	client := bitbucket.NewOAuthbearerToken(token)
	diffOptions := &bitbucket.DiffOptions{
		Owner:    repo.Namespace,
		RepoSlug: repo.Name,
		Spec:     fmt.Sprintf("%s...%s", build.Before, build.After),
	}
	diffs, err := client.Repositories.Diff.GetDiff(diffOptions)
	if err != nil {
		return nil, err
	}
	logrus.Info(diffs)
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
