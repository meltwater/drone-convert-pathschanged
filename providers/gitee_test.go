package providers

import (
	"github.com/h2non/gock"
	"testing"

	"github.com/drone/drone-go/drone"
	"github.com/drone/drone-go/plugin/converter"

	"github.com/google/go-cmp/cmp"
)

func TestGetGiteeFilesChangedCommit(t *testing.T) {
	defer gock.Off()

	gock.New("https://gitee.com/api/v5").
		Get("/repos/kit101/drone-yml-test/commits/537575f44a09c57dfc472e26fe067754fd2f9374").
		Reply(200).
		Type("application/json").
		File("testdata/gitee/commit.json")

	req := &converter.Request{
		Build: drone.Build{
			Before: "",
			After:  "537575f44a09c57dfc472e26fe067754fd2f9374",
		},
		Repo: drone.Repo{
			Namespace: "kit101",
			Name:      "drone-yml-test",
			Slug:      "kit101/drone-yml-test",
			Config:    ".drone.yml",
		},
	}

	got, err := GetGiteeFilesChanged(req.Repo, req.Build, "invalidtoken")
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

func TestGetGiteeFilesChangedCompare(t *testing.T) {
	defer gock.Off()

	gock.New("https://gitee.com/api/v5").
		Get("/repos/kit101/drone-yml-test/commits/e3c0ff4d5cef439ea11b30866fb1ed79b420801d").
		Reply(200).
		Type("application/json").
		File("testdata/gitee/compare.json")

	req := &converter.Request{
		Build: drone.Build{
			Before: "",
			After:  "e3c0ff4d5cef439ea11b30866fb1ed79b420801d",
		},
		Repo: drone.Repo{
			Namespace: "kit101",
			Name:      "drone-yml-test",
			Slug:      "kit101/drone-yml-test",
			Config:    ".drone.yml",
		},
	}

	got, err := GetGiteeFilesChanged(req.Repo, req.Build, "invalidtoken")
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
