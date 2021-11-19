package providers

import (
	"context"
	"net/http"

	"github.com/drone/drone-go/drone"
	"github.com/drone/go-scm/scm"
	"github.com/drone/go-scm/scm/driver/bitbucket"
	"github.com/drone/go-scm/scm/transport"
)

func GetBitbucketFilesChanged(repo drone.Repo, build drone.Build, user string, password string, opts scm.ListOptions) ([]string, error) {
	newctx := context.Background()

	client, err := bitbucket.New("https://api.bitbucket.org")
	if err != nil {
		return nil, err
	}

	client.Client = &http.Client{
		Transport: &transport.BasicAuth{
			Username: user,
			Password: password,
		},
	}

	var got []*scm.Change

	if build.Before == "" || build.Before == scm.EmptyCommit {
		got, _, err = client.Git.ListChanges(newctx, repo.Slug, build.After, opts)
		if err != nil {
			return nil, err
		}
	} else {
		got, _, err = client.Git.CompareChanges(newctx, repo.Slug, build.Before, build.After, opts)
		if err != nil {
			return nil, err
		}
	}

	var files []string
	for _, c := range got {
		files = append(files, c.Path)
	}

	return files, nil
}
