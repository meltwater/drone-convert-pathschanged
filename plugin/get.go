package plugin

import (
	"context"

	"github.com/drone/drone-go/drone"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

func getFilesChanged(repo drone.Repo, build drone.Build, token string) ([]string, error) {
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

	var files []string
	for _, f := range commitFiles {
		files = append(files, *f.Filename)
	}

	return files, nil
}

func apiRateLimit () {
    newctx := context.Background()
    ts := oauth2.StaticTokenSource(
	&oauth2.Token{AccessToken: token},
    )
    tc := oauth2.NewClient(newctx, ts)

	client := github.NewClient(tc)
	
	rateLimit, _, err := client.RateLimits(newctx)
	if err != nil {
		fmt.Printf("Promblem getting github rate limit info %v\n", err)
		return
	}
	GithubApiCount = (rateLimit.Core.Limit - rateLimit.Core.Remaining)
    return GithubApiCount, nil
}
