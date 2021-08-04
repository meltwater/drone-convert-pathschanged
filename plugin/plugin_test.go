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

	providers := []string{"github", "bitbucket-server", "gitea"}

	req := &converter.Request{
		Build: drone.Build{},
		Repo: drone.Repo{
			Slug:   "somewhere/over-the-rainbow",
			Config: ".drone.yml",
		},
	}

	for _, provider := range providers {
		plugin := New("invalidtoken", provider, "", "")

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

	plugin := New("invalidtoken", "", "", "")

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

	plugin := New("invalidtoken", "unsupported", "", "")

	_, err := plugin.Convert(noContext, req)
	if err == nil {
		t.Error("unsupported provider did not return error")
		return
	}
}

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

	plugin := New("invalidtoken", "github", "", "")

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

	plugin := New("invalidtoken", "github", "", "")

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

	plugin := New("invalidtoken", "github", "", "")

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

	plugin := New("invalidtoken", "github", "", "")

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