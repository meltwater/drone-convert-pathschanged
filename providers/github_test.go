package providers

import (
	"testing"

	"github.com/drone/drone-go/drone"
	"github.com/drone/drone-go/plugin/converter"

	"github.com/google/go-cmp/cmp"
	"github.com/h2non/gock"
)

var mockHeaders = map[string]string{
	"X-GitHub-Request-Id":   "DD0E:6011:12F21A8:1926790:5A2064E2",
	"X-RateLimit-Limit":     "60",
	"X-RateLimit-Remaining": "59",
	"X-RateLimit-Reset":     "1512076018",
}

func TestGetGithubFilesChangedCommit(t *testing.T) {
	defer gock.Off()

	gock.New("https://api.github.com").
		Get("/repos/meltwater/drone-convert-pathschanged/commits/6ee3cf41d995a79857e0db41c47bf619e6546571").
		Reply(200).
		Type("application/json").
		SetHeaders(mockHeaders).
		File("testdata/github/commit.json")

	req := &converter.Request{
		Build: drone.Build{
			Before: "",
			After:  "6ee3cf41d995a79857e0db41c47bf619e6546571",
		},
		Repo: drone.Repo{
			Namespace: "meltwater",
			Name:      "drone-convert-pathschanged",
			Slug:      "meltwater/drone-convert-pathschanged",
			Config:    ".drone.yml",
		},
	}

	got, err := GetGithubFilesChanged(req.Repo, req.Build, "invalidtoken")
	if err != nil {
		t.Error(err)
		return
	}

	want := []string{".drone.yml"}

	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("Unexpected Results")
		t.Log(diff)
	}
}

func TestGetGithubFilesChangedCompare(t *testing.T) {
	defer gock.Off()

	gock.New("https://api.github.com").
		Get("/repos/meltwater/drone-convert-pathschanged/compare/496eb80334e84085426ce681407d770cc9247acd...6ee3cf41d995a79857e0db41c47bf619e6546571").
		Reply(200).
		Type("application/json").
		SetHeaders(mockHeaders).
		File("testdata/github/compare.json")

	req := &converter.Request{
		Build: drone.Build{
			Before: "496eb80334e84085426ce681407d770cc9247acd",
			After:  "6ee3cf41d995a79857e0db41c47bf619e6546571",
		},
		Repo: drone.Repo{
			Namespace: "meltwater",
			Name:      "drone-convert-pathschanged",
			Slug:      "meltwater/drone-convert-pathschanged",
			Config:    ".drone.yml",
		},
	}

	got, err := GetGithubFilesChanged(req.Repo, req.Build, "invalidtoken")
	if err != nil {
		t.Error(err)
		return
	}

	want := []string{".drone.yml"}

	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("Unexpected Results")
		t.Log(diff)
	}
}
