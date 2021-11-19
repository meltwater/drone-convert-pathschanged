package providers

import (
	"context"
	"net/http"

	"github.com/drone/drone-go/drone"
	"github.com/drone/go-scm/scm"
	"github.com/drone/go-scm/scm/driver/bitbucket"
	"github.com/drone/go-scm/scm/transport"
	"github.com/sirupsen/logrus"
)

func GetBitbucketFilesChanged(repo drone.Repo, build drone.Build, user string, password string) ([]string, error) {
	newctx := context.Background()

	client, err := bitbucket.New("https://api.bitbucket.org")
	if err != nil {
		return nil, err
	}

	requestLogger := logrus.WithFields(logrus.Fields{
		"repo_namespace": repo.Namespace,
		"repo_name":      repo.Name,
	})

	client.Client = &http.Client{
		Transport: &transport.BasicAuth{
			Username: user,
			Password: password,
		},
	}

	//got, _, err := client.Git.ListChanges(newctx, repo.Slug, build.After, scm.ListOptions{})
	got, _, err := client.Git.CompareChanges(newctx, repo.Slug, build.Before, build.After, scm.ListOptions{})
	if err != nil {
		return nil, err
	}

	var files []string
	for _, c := range got {
		files = append(files, c.Path)
	}
	requestLogger.Infoln("files", files)

	return files, nil
}
