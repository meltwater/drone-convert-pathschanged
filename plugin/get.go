package plugin

import (
	"context"

	"github.com/drone/drone-go/drone"
	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/sirupsen/logrus"
)

func getFilesChanged(repo drone.Repo, build drone.Build, token string, server string) ([]string, error) {
	newctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	trans := oauth2.NewClient(newctx, ts)

	var client *github.Client

	if server == "" {
		client = github.NewClient(trans)
	} else {
		var err error
		client, err = github.NewEnterpriseClient(server, server, trans)
		if err != nil {
			logrus.Errorf("Unable to connect to Github: '%v'", err)
			return nil, err
		}
	}

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
    rateLimit, _,err := client.RateLimits(newctx)
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

var (
	GithubApiCount = prometheus.NewGauge(
	prometheus.GaugeOpts{
		   Name: "github_api_calls_remaining",
		   Help: "Total number of github api calls per hour remaining",
})
)
func init(){
    prometheus.MustRegister(GithubApiCount)
}
