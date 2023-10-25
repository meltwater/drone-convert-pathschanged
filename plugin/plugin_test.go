package plugin

import (
	"context"
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

// empty context
var noContext = context.Background()

func TestNewEmptyPipeline(t *testing.T) {

	providers := []string{"github", "bitbucket", "bitbucket-server", "stash"}

	req := &converter.Request{
		Build: drone.Build{},
		Repo: drone.Repo{
			Slug:   "somewhere/over-the-rainbow",
			Config: ".drone.yml",
		},
	}

	params := &Params{
		Token: "invalidtoken",
	}

	for _, provider := range providers {
		plugin := New(provider, params)

		config, err := plugin.Convert(noContext, req)
		if err != nil {
			t.Error(err)
			return
		}

		if want, got := "", config.Data; want != got {
			t.Errorf("Want %q got %q", want, got)
		}
	}
}

func TestNewInvalidPipeline(t *testing.T) {
	data := `
kind: pipeline
type: docker
name: default

this_is_invalid_yaml
`

	req := &converter.Request{
		Build: drone.Build{
			Before: "",
			After:  "6ee3cf41d995a79857e0db41c47bf619e6546571",
		},
		Config: drone.Config{
			Data: data,
		},
		Repo: drone.Repo{
			Namespace: "meltwater",
			Name:      "drone-convert-pathschanged",
			Slug:      "meltwater/drone-convert-pathschanged",
			Config:    ".drone.yml",
		},
	}

	params := &Params{
		Token: "invalidtoken",
	}

	plugin := New("invalidtoken", params)

	_, err := plugin.Convert(noContext, req)
	if err == nil {
		t.Error("invalid pipeline did not return error")
		return
	}
}

func TestNewUnsupportedProvider(t *testing.T) {

	data := `
kind: pipeline
type: docker
name: default

steps:
- name: message
  image: busybox
  commands:
  - echo "This step will be excluded when .drone.yml is changed"
  when:
    paths:
      exclude:
      - .drone.yml
`

	params := &Params{
		Token: "invalidtoken",
	}

	req := &converter.Request{
		Build: drone.Build{
			Before: "",
			After:  "6ee3cf41d995a79857e0db41c47bf619e6546571",
		},
		Config: drone.Config{
			Data: data,
		},
		Repo: drone.Repo{
			Namespace: "meltwater",
			Name:      "drone-convert-pathschanged",
			Slug:      "meltwater/drone-convert-pathschanged",
			Config:    ".drone.yml",
		},
	}

	plugin := New("unsupported", params)

	_, err := plugin.Convert(noContext, req)

	// just looking for an error here isn't enough, since calling 'New' will always return
	// an error during this test, because it can't authenticate with the provider
	//
	// therefore, look for the specific 'unsupported provided' error
	if err.Error() != "unsupported provider" {
		t.Error("unsupported provider did not return error")
		return
	}
}

// github provider
func TestNewGithubCommitExcludeStep(t *testing.T) {
	gock.New("https://api.github.com").
		Get("/repos/meltwater/drone-convert-pathschanged/commits/6ee3cf41d995a79857e0db41c47bf619e6546571").
		Reply(200).
		Type("application/json").
		SetHeaders(mockHeaders).
		File("../providers/testdata/github/commit.json")

	before := `
kind: pipeline
type: docker
name: default

steps:
- name: message
  image: busybox
  commands:
  - echo "This step will be excluded when .drone.yml is changed"
  when:
    paths:
      exclude:
      - .drone.yml
`

	params := &Params{
		Token: "invalidtoken",
	}

	req := &converter.Request{
		Build: drone.Build{
			Before: "",
			After:  "6ee3cf41d995a79857e0db41c47bf619e6546571",
		},
		Config: drone.Config{
			Data: before,
		},
		Repo: drone.Repo{
			Namespace: "meltwater",
			Name:      "drone-convert-pathschanged",
			Slug:      "meltwater/drone-convert-pathschanged",
			Config:    ".drone.yml",
		},
	}

	plugin := New("github", params)

	config, err := plugin.Convert(noContext, req)
	if err != nil {
		t.Error(err)
		return
	}

	after := `kind: pipeline
type: docker
steps:
- when:
    paths:
      exclude:
      - .drone.yml
    event:
      exclude:
      - '*'
  commands:
  - echo "This step will be excluded when .drone.yml is changed"
  image: busybox
  name: message
name: default
`
	want := &drone.Config{
		Data: after,
	}

	if diff := cmp.Diff(config.Data, want.Data); diff != "" {
		t.Errorf("Unexpected Results")
		t.Log(diff)
	}
}

func TestNewGithubCommitIncludeStep(t *testing.T) {
	gock.New("https://api.github.com").
		Get("/repos/meltwater/drone-convert-pathschanged/commits/6ee3cf41d995a79857e0db41c47bf619e6546571").
		Reply(200).
		Type("application/json").
		SetHeaders(mockHeaders).
		File("../providers/testdata/github/commit.json")

	before := `
kind: pipeline
type: docker
name: default

steps:
- name: message
  image: busybox
  commands:
  - echo "This step will be included when .drone.yml is changed"
  when:
    paths:
      include:
      - .drone.yml
`

	params := &Params{
		Token: "invalidtoken",
	}

	req := &converter.Request{
		Build: drone.Build{
			Before: "",
			After:  "6ee3cf41d995a79857e0db41c47bf619e6546571",
		},
		Config: drone.Config{
			Data: before,
		},
		Repo: drone.Repo{
			Namespace: "meltwater",
			Name:      "drone-convert-pathschanged",
			Slug:      "meltwater/drone-convert-pathschanged",
			Config:    ".drone.yml",
		},
	}

	plugin := New("github", params)

	config, err := plugin.Convert(noContext, req)
	if err != nil {
		t.Error(err)
		return
	}

	after := `kind: pipeline
type: docker
steps:
- when:
    paths:
      include:
      - .drone.yml
  commands:
  - echo "This step will be included when .drone.yml is changed"
  image: busybox
  name: message
name: default
`
	want := &drone.Config{
		Data: after,
	}

	if diff := cmp.Diff(config.Data, want.Data); diff != "" {
		t.Errorf("Unexpected Results")
		t.Log(diff)
	}
}

func TestNewGithubCompareExcludeStep(t *testing.T) {
	gock.New("https://api.github.com").
		Get("/repos/meltwater/drone-convert-pathschanged/compare/496eb80334e84085426ce681407d770cc9247acd...6ee3cf41d995a79857e0db41c47bf619e6546571").
		Reply(200).
		Type("application/json").
		SetHeaders(mockHeaders).
		File("../providers/testdata/github/compare.json")

	before := `
kind: pipeline
type: docker
name: default

steps:
- name: message
  image: busybox
  commands:
  - echo "This step will be excluded when .drone.yml is changed"
  when:
    paths:
      exclude:
      - .drone.yml
`

	params := &Params{
		Token: "invalidtoken",
	}

	req := &converter.Request{
		Build: drone.Build{
			Before: "496eb80334e84085426ce681407d770cc9247acd",
			After:  "6ee3cf41d995a79857e0db41c47bf619e6546571",
		},
		Config: drone.Config{
			Data: before,
		},
		Repo: drone.Repo{
			Namespace: "meltwater",
			Name:      "drone-convert-pathschanged",
			Slug:      "meltwater/drone-convert-pathschanged",
			Config:    ".drone.yml",
		},
	}

	plugin := New("github", params)

	config, err := plugin.Convert(noContext, req)
	if err != nil {
		t.Error(err)
		return
	}

	after := `kind: pipeline
type: docker
steps:
- when:
    paths:
      exclude:
      - .drone.yml
    event:
      exclude:
      - '*'
  commands:
  - echo "This step will be excluded when .drone.yml is changed"
  image: busybox
  name: message
name: default
`
	want := &drone.Config{
		Data: after,
	}

	if diff := cmp.Diff(config.Data, want.Data); diff != "" {
		t.Errorf("Unexpected Results")
		t.Log(diff)
	}
}

func TestNewGithubCompareIncludeStep(t *testing.T) {
	gock.New("https://api.github.com").
		Get("/repos/meltwater/drone-convert-pathschanged/compare/496eb80334e84085426ce681407d770cc9247acd...6ee3cf41d995a79857e0db41c47bf619e6546571").
		Reply(200).
		Type("application/json").
		SetHeaders(mockHeaders).
		File("../providers/testdata/github/compare.json")

	before := `
kind: pipeline
type: docker
name: default

steps:
- name: message
  image: busybox
  commands:
  - echo "This step will be included when .drone.yml is changed"
  when:
    paths:
      include:
      - .drone.yml
`
	params := &Params{
		Token: "invalidtoken",
	}

	req := &converter.Request{
		Build: drone.Build{
			Before: "496eb80334e84085426ce681407d770cc9247acd",
			After:  "6ee3cf41d995a79857e0db41c47bf619e6546571",
		},
		Config: drone.Config{
			Data: before,
		},
		Repo: drone.Repo{
			Namespace: "meltwater",
			Name:      "drone-convert-pathschanged",
			Slug:      "meltwater/drone-convert-pathschanged",
			Config:    ".drone.yml",
		},
	}

	plugin := New("github", params)

	config, err := plugin.Convert(noContext, req)
	if err != nil {
		t.Error(err)
		return
	}

	after := `kind: pipeline
type: docker
steps:
- when:
    paths:
      include:
      - .drone.yml
  commands:
  - echo "This step will be included when .drone.yml is changed"
  image: busybox
  name: message
name: default
`
	want := &drone.Config{
		Data: after,
	}

	if diff := cmp.Diff(config.Data, want.Data); diff != "" {
		t.Errorf("Unexpected Results")
		t.Log(diff)
	}
}

// bitbucket provider
func TestNewBitbucketCommitExcludeStep(t *testing.T) {
	gock.New("https://api.bitbucket.org").
		Get("/2.0/repositories/atlassian/atlaskit/diffstat/425863f9dbe56d70c8dcdbf2e4e0805e85591fcc").
		Reply(200).
		Type("application/json").
		File("../providers/testdata/bitbucket/diffstat.json")

	before := `
kind: pipeline
type: docker
name: default

steps:
- name: message
  image: busybox
  commands:
  - echo "This step will be excluded when CONTRIBUTING.md is changed"
  when:
    paths:
      exclude:
      - CONTRIBUTING.md
`

	params := &Params{
		BitBucketUser:     "centauri",
		BitBucketPassword: "kodan",
		Token:             "invalidtoken",
	}

	req := &converter.Request{
		Build: drone.Build{
			Before: "",
			After:  "425863f9dbe56d70c8dcdbf2e4e0805e85591fcc",
		},
		Config: drone.Config{
			Data: before,
		},
		Repo: drone.Repo{
			Namespace: "atlassian",
			Name:      "atlaskit",
			Slug:      "atlassian/atlaskit",
			Config:    ".drone.yml",
		},
	}

	plugin := New("bitbucket", params)

	config, err := plugin.Convert(noContext, req)
	if err != nil {
		t.Error(err)
		return
	}

	after := `kind: pipeline
type: docker
steps:
- when:
    paths:
      exclude:
      - CONTRIBUTING.md
    event:
      exclude:
      - '*'
  commands:
  - echo "This step will be excluded when CONTRIBUTING.md is changed"
  image: busybox
  name: message
name: default
`
	want := &drone.Config{
		Data: after,
	}

	if diff := cmp.Diff(config.Data, want.Data); diff != "" {
		t.Errorf("Unexpected Results")
		t.Log(diff)
	}
}

func TestNewBitbucketCompareExcludeStep(t *testing.T) {
	gock.New("https://api.bitbucket.org").
		Get("/2.0/repositories/atlassian/atlaskit/diffstat/425863f9dbe56d70c8dcdbf2e4e0805e85591fcc..dec26e0fe887167743c2b7e36531dedfeb6cd478").
		Reply(200).
		Type("application/json").
		File("../providers/testdata/bitbucket/diffstat.json")

	before := `
kind: pipeline
type: docker
name: default

steps:
- name: message
  image: busybox
  commands:
  - echo "This step will be excluded when CONTRIBUTING.md is changed"
  when:
    paths:
      exclude:
      - CONTRIBUTING.md
`
	params := &Params{
		BitBucketUser:     "centauri",
		BitBucketPassword: "kodan",
		Token:             "invalidtoken",
	}

	req := &converter.Request{
		Build: drone.Build{
			Before: "dec26e0fe887167743c2b7e36531dedfeb6cd478",
			After:  "425863f9dbe56d70c8dcdbf2e4e0805e85591fcc",
		},
		Config: drone.Config{
			Data: before,
		},
		Repo: drone.Repo{
			Namespace: "atlassian",
			Name:      "atlaskit",
			Slug:      "atlassian/atlaskit",
			Config:    ".drone.yml",
		},
	}

	plugin := New("bitbucket", params)

	config, err := plugin.Convert(noContext, req)
	if err != nil {
		t.Error(err)
		return
	}

	after := `kind: pipeline
type: docker
steps:
- when:
    paths:
      exclude:
      - CONTRIBUTING.md
    event:
      exclude:
      - '*'
  commands:
  - echo "This step will be excluded when CONTRIBUTING.md is changed"
  image: busybox
  name: message
name: default
`
	want := &drone.Config{
		Data: after,
	}

	if diff := cmp.Diff(config.Data, want.Data); diff != "" {
		t.Errorf("Unexpected Results")
		t.Log(diff)
	}
}

func TestStashPageSizeSetWhenConfigured(t *testing.T) {
	gock.New("http://example.com:7990").
		Get("/rest/api/1.0/projects/PRJ/repos/my-repo/compare/changes").
		MatchParam("from", "4f4b0ef1714a5b6cafdaf2f53c7f5f5b38fb9348").
		MatchParam("to", "131cb13f4aed12e725177bc4b7c28db67839bf9f").
		MatchParam("limit", "40"). //this will resolve only if the configured page size was honoured
		Reply(200).
		Type("application/json").
		File("../providers/testdata/stash/compare.json")

	before := `
kind: pipeline
type: docker
name: default

steps:
- name: message
  image: busybox
  commands:
  - echo "This step will be excluded when .drone.yml is changed"
  when:
    paths:
      exclude:
      - .drone.yml
`
	params := &Params{
		StashServer:   "http://example.com:7990",
		Token:         "invalidtoken",
		StashPageSize: 40,
	}

	req := &converter.Request{
		Build: drone.Build{
			Before: "4f4b0ef1714a5b6cafdaf2f53c7f5f5b38fb9348",
			After:  "131cb13f4aed12e725177bc4b7c28db67839bf9f",
		},
		Config: drone.Config{
			Data: before,
		},
		Repo: drone.Repo{
			Namespace: "repos",
			Name:      "my-repo",
			Slug:      "PRJ/my-repo",
			Config:    ".drone.yml",
		},
	}

	plugin := New("stash", params)

	config, err := plugin.Convert(noContext, req)
	if err != nil {
		t.Error(err)
		return
	}

	if config == nil { //we are not interested in transformation, only validating limit was set on the request to stash and resolved to gock
		t.Errorf("Unexpected Results")
	}

}
