package providers

import (
	"context"
	"encoding/json"

	"github.com/drone/drone-go/drone"
	bitbucketv1 "github.com/gfleury/go-bitbucket-v1"
)

type bitbucketDiffs struct {
	Diffs []struct {
		Destination struct {
			Name string `json:"toString"`
		} `json:"destination"`
	} `json:"diffs"`
}

func GetBBFilesChanged(repo drone.Repo, build drone.Build, token string, bitbucketAddress string) ([]string, error) {
	var files []string
	var ctx context.Context
	params := map[string]interface{}{
		"since": build.Before,
	}
	ctx = context.WithValue(context.Background(), bitbucketv1.ContextAccessToken, token)

	configuration := bitbucketv1.NewConfiguration(bitbucketAddress)
	client := bitbucketv1.NewAPIClient(ctx, configuration)

	ff, err := client.DefaultApi.StreamDiff(repo.Namespace, repo.Name, build.After, params)
	if err != nil {
		return nil, err
	}

	jsonString, err := json.Marshal(ff.Values)
	if err != nil {
		return nil, err
	}
	res := bitbucketDiffs{}
	err = json.Unmarshal([]byte(jsonString), &res)
	if err != nil {
		return nil, err
	}

	for _, diff := range res.Diffs {
		files = append(files, diff.Destination.Name)
	}
	return files, nil
}