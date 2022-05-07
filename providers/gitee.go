// create by vinson on 2020/4/10

package providers

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/drone/drone-go/drone"
	"github.com/drone/go-scm/scm"
	"github.com/drone/go-scm/scm/driver/gitee"
	"github.com/drone/go-scm/scm/transport"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

func GetGiteeFilesChanged(repo drone.Repo, build drone.Build, username, password, clientID, clientSecret, scope string) ([]string, error) {
	newctx := context.Background()
	var client *scm.Client
	var err error

	client = gitee.NewDefault()

	postData := url.Values{}
	postData.Add("grant_type", "password")
	postData.Add("username", strings.TrimSpace(username))
	postData.Add("password", strings.TrimSpace(password))
	postData.Add("client_id", strings.TrimSpace(clientID))
	postData.Add("client_secret", strings.TrimSpace(clientSecret))
	postData.Add("scope", strings.TrimSpace(scope))

	res, err := http.Post("https://gitee.com/oauth/token", "application/x-www-form-urlencoded", strings.NewReader(postData.Encode()))
	if res == nil {
		fmt.Println("Fatal error : res is nil")
		return nil, err
	}
	if err != nil {
		fmt.Println("Fatal error ", err.Error())
		return nil, err
	}

	if res.StatusCode != 200 {
		fmt.Println("Fatal error ", res.StatusCode)
		return nil, err
	}

	defer func(Body io.ReadCloser) {
		e := Body.Close()
		if e != nil {
			fmt.Println("Fatal error ", err.Error())
		}
	}(res.Body)

	content, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println("Fatal error ", err.Error())
	}

	var str = make(map[string]interface{}, 0)
	err = json.Unmarshal(content, &str)
	if err != nil {
		fmt.Println("Fatal error ", err.Error())
	}

	client.Client = &http.Client{
		Transport: &transport.BearerToken{
			Token: str["access_token"].(string),
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
