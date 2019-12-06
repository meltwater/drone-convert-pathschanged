package plugin

import (
	"context"
	"regexp"
	"strconv"

	"github.com/drone/drone-go/drone"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

// regular expression to extract the pull request number
// from the git ref (e.g. refs/pulls/{d}/head)
var re = regexp.MustCompile("\\d+")

func getFilesChanged(repo drone.Repo, build drone.Build, token string) ([]string, error) {
	newctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(newctx, ts)

	client := github.NewClient(tc)

	var commitFiles []github.CommitFile

	if build.Event == "pull_request" {
		s := re.FindString(build.Ref)
		d, _ := strconv.Atoi(s)

		listOpts := &github.ListOptions{
			Page:    1,
			PerPage: 300,
		}
		allFiles := []*github.CommitFile{}

		for {
			cfs, resp, err := client.PullRequests.ListFiles(newctx, repo.Namespace, repo.Name, d, listOpts)
			if err != nil {
				return nil, err
			}

			allFiles = append(allFiles, cfs...)
			if resp.NextPage == 0 {
				break
			}
			listOpts.Page = resp.NextPage
		}
		for _, file := range allFiles {
			commitFiles = append(commitFiles, *file)
		}
	} else if build.Before == "" || build.Before == "0000000000000000000000000000000000000000" {
		rc, _, err := client.Repositories.GetCommit(newctx, repo.Namespace, repo.Name, build.After)
		if err != nil {
			return nil, err
		}
		commitFiles = rc.Files
	} else {
		cc, _, err := client.Repositories.CompareCommits(newctx, repo.Namespace, repo.Name, build.Before, build.After)
		if err != nil {
			return nil, err
		}
		commitFiles = cc.Files
	}

	var files []string
	for _, f := range commitFiles {
		files = append(files, *f.Filename)
	}

	return files, nil
}
