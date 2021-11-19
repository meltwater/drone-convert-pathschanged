package providers

import (
	"testing"

	"github.com/drone/drone-go/drone"
	"github.com/drone/drone-go/plugin/converter"
	"github.com/drone/go-scm/scm"

	"github.com/google/go-cmp/cmp"
	"github.com/h2non/gock"
)

func TestGetBitbucketFilesChangedCommit(t *testing.T) {
	defer gock.Off()

	gock.New("https://api.bitbucket.org").
		Get("/2.0/repositories/atlassian/atlaskit/diffstat/425863f9dbe56d70c8dcdbf2e4e0805e85591fcc").
		MatchParam("page", "1").
		MatchParam("pagelen", "30").
		Reply(200).
		Type("application/json").
		File("testdata/bitbucket/diffstat.json")

	req := &converter.Request{
		Build: drone.Build{
			Before: "",
			After:  "425863f9dbe56d70c8dcdbf2e4e0805e85591fcc",
		},
		Repo: drone.Repo{
			Namespace: "atlassian",
			Name:      "atlaskit",
			Slug:      "atlassian/atlaskit",
			Config:    ".drone.yml",
		},
	}

	got, err := GetBitbucketFilesChanged(req.Repo, req.Build, "centauri", "kodan", scm.ListOptions{Page: 1, Size: 30})
	if err != nil {
		t.Error(err)
		return
	}

	want := []string{"CONTRIBUTING.md"}

	if diff := cmp.Diff(got, want); diff != "" {
		t.Error("Unexpected Results")
		t.Log(diff)
	}
}

func TestGetBitbucketFilesChangedCompare(t *testing.T) {
	defer gock.Off()

	gock.New("https://api.bitbucket.org").
		Get("/2.0/repositories/atlassian/atlaskit/diffstat/dec26e0fe887167743c2b7e36531dedfeb6cd478..425863f9dbe56d70c8dcdbf2e4e0805e85591fcc").
		MatchParam("page", "1").
		MatchParam("pagelen", "30").
		Reply(200).
		Type("application/json").
		File("testdata/bitbucket/diffstat.json")

	// build.Before and build.After are switched due to a bug https://github.com/drone/go-scm/pull/127
	// FIXME: switcch build.Before and build.After parameters when the above issue is fixed
	req := &converter.Request{
		Build: drone.Build{
			Before: "425863f9dbe56d70c8dcdbf2e4e0805e85591fcc",
			After:  "dec26e0fe887167743c2b7e36531dedfeb6cd478",
		},
		Repo: drone.Repo{
			Namespace: "atlassian",
			Name:      "atlaskit",
			Slug:      "atlassian/atlaskit",
			Config:    ".drone.yml",
		},
	}

	got, err := GetBitbucketFilesChanged(req.Repo, req.Build, "centauri", "kodan", scm.ListOptions{Page: 1, Size: 30})
	if err != nil {
		t.Error(err)
		return
	}

	want := []string{"CONTRIBUTING.md"}

	if diff := cmp.Diff(got, want); diff != "" {
		t.Error("Unexpected Results")
		t.Log(diff)
	}
}
