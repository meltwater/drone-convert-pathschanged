package plugin

import (
	"fmt"
	"github.com/drone/drone-go/drone"
	_ "github.com/ktrysmt/go-bitbucket"
	"github.com/meltwater/drone-convert-pathschanged/providers"
)

func getFilesChanged(repo drone.Repo, build drone.Build, token string, provider string) ([]string, error) {
	switch provider {
	case "github":
		changedFiles, err := providers.GetGithubFilesChanged(repo, build, token)
		if err != nil {
			return nil, err
		}
		return changedFiles, nil
	case "bitbucket-server":
		changedFiles, err := providers.GetBBFilesChanged(repo, build, token)
		if err != nil {
			return nil, err
		}
		return changedFiles, nil
	default:
		return nil, fmt.Errorf("unsupported provider %s", provider)
	}
}
