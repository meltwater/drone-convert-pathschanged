package providers

import (
	"testing"

	"github.com/drone/drone-go/drone"
	"github.com/drone/drone-go/plugin/converter"
	"github.com/drone/go-scm/scm"
	"github.com/google/go-cmp/cmp"
	"github.com/h2non/gock"
)

func TestGetStashFilesChangedCommit(t *testing.T) {
	defer gock.Off()

	gock.New("http://example.com:7990").
		Get("/rest/api/1.0/projects/PRJ/repos/my-repo/commits/131cb13f4aed12e725177bc4b7c28db67839bf9f/changes").
		MatchParam("limit", "30").
		Reply(200).
		Type("application/json").
		File("testdata/stash/changes.json")

	req := &converter.Request{
		Build: drone.Build{
			Before: "",
			After:  "131cb13f4aed12e725177bc4b7c28db67839bf9f",
		},
		Repo: drone.Repo{
			Namespace: "repos",
			Name:      "my-repo",
			Slug:      "PRJ/my-repo",
			Config:    ".drone.yml",
		},
	}

	got, err := GetStashFilesChanged(req.Repo, req.Build, "http://example.com:7990", "invalidtoken", scm.ListOptions{Page: 1, Size: 30})
	if err != nil {
		t.Error(err)
		return
	}

	want := []string{
		".gitignore",
		"COPYING",
		"README.md",
		"main.go",
	}

	if diff := cmp.Diff(got, want); diff != "" {
		t.Error("Unexpected Results")
		t.Log(diff)
	}
}
