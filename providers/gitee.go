// create by vinson on 2020/4/10

package providers

import (
	"context"
	"github.com/drone/drone-go/drone"
	"github.com/drone/go-scm/scm"
	"github.com/drone/go-scm/scm/driver/gitee"
	"github.com/drone/go-scm/scm/transport"
	"net/http"
)

func GetGiteeFilesChanged(repo drone.Repo, build drone.Build, token string) ([]string, error) {
	newctx := context.Background()
	var client *scm.Client
	var err error

	client = gitee.NewDefault()

	client.Client = &http.Client{
		Transport: &transport.BearerToken{
			Token: token,
		},
	}

	var changes []*scm.Change
	var _ *scm.Response

	if build.Before == "" || build.Before == scm.EmptyCommit {
		changes, _, err = client.Git.ListChanges(newctx, repo.Slug, build.After, scm.ListOptions{})
		if err != nil {
			return nil, err
		}
	} else {
		changes, _, err = client.Git.CompareChanges(newctx, repo.Slug, build.Before, build.After, scm.ListOptions{})
		if err != nil {
			return nil, err
		}
	}

	var files []string
	for _, c := range changes {
		files = append(files, c.Path)
	}

	return files, nil
}
