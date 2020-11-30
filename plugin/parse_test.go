package plugin

import (
	"testing"

	"github.com/drone/drone-go/drone"
	"github.com/drone/drone-go/plugin/converter"
)

// pathSeen tests
func TestPathSeenEmptyPipeline(t *testing.T) {

	data := `
kind: pipeline
type: docker
name: default
`

	pathSeen, err := pathSeen(data)
	if err != nil {
		t.Error(err)
		return
	}

	if want, got := false, pathSeen; want != got {
		t.Errorf("Want %t got %t", want, got)
	}
}

func TestPathSeenInvalidPipeline(t *testing.T) {

	data := `
kind: pipeline
type: docker
name: default

this_is_invalid_yaml
`

	_, err := pathSeen(data)

	if err == nil {
		t.Errorf("Invalid pipeline did not return an error")
	}
}

func TestPathSeenSingleStepPipeline(t *testing.T) {

	data := `
kind: pipeline
type: docker
name: default

trigger:
  branch:
  - master

steps:
- name: message
  image: busybox
  commands:
  - echo "hello"
`

	pathSeen, err := pathSeen(data)
	if err != nil {
		t.Error(err)
		return
	}

	if want, got := false, pathSeen; want != got {
		t.Errorf("Want %t got %t", want, got)
	}
}

func TestPathSeenStepPathIncludePipeline(t *testing.T) {

	data := `
kind: pipeline
type: docker
name: default

steps:
- name: message
  image: busybox
  commands:
  - echo "README.md was changed"
  when:
    paths:
      include:
      - README.md
`

	pathSeen, err := pathSeen(data)
	if err != nil {
		t.Error(err)
		return
	}

	if want, got := true, pathSeen; want != got {
		t.Errorf("Want %t got %t", want, got)
	}
}

func TestPathSeenStepPathImplicitAndOptionalIncludePipeline(t *testing.T) {

	data := `
kind: pipeline
type: docker
name: default

steps:
- name: message
  image: busybox
  commands:
  - echo "README.md was changed"
  when:
    paths:
      - README.md
`

	pathSeen, err := pathSeen(data)
	if err != nil {
		t.Error(err)
		return
	}

	if !pathSeen {
		t.Errorf("Want 'true' got %t", pathSeen)
	}
}

func TestPathSeenStepPathExcludePipeline(t *testing.T) {

	data := `
kind: pipeline
type: docker
name: default

trigger:
  paths:
    exclude:
    - README.md

steps:
- name: message
  image: busybox
  commands:
  - echo "This pipeline will be excluded when README.md is changed"
`

	pathSeen, err := pathSeen(data)
	if err != nil {
		t.Error(err)
		return
	}

	if want, got := true, pathSeen; want != got {
		t.Errorf("Want %t got %t", want, got)
	}
}

func TestPathSeenStepPathIncludeAndExcludePipeline(t *testing.T) {

	data := `
kind: pipeline
type: docker
name: default

steps:
- name: message
  image: busybox
  commands:
  - echo "CHANGELOG.md was changed, README.md was not changed"
  when:
    paths:
      include:
      - CHANGELOG.md
      exclude:
      - README.md
`

	pathSeen, err := pathSeen(data)
	if err != nil {
		t.Error(err)
		return
	}

	if want, got := true, pathSeen; want != got {
		t.Errorf("Want %t got %t", want, got)
	}
}

func TestPathSeenTriggerPathIncludePipeline(t *testing.T) {

	data := `
kind: pipeline
type: docker
name: default

trigger:
  paths:
    include:
    - README.md

steps:
- name: message
  image: busybox
  commands:
  - echo "README.md was changed"
`
	pathSeen, err := pathSeen(data)
	if err != nil {
		t.Error(err)
		return
	}

	if want, got := true, pathSeen; want != got {
		t.Errorf("Want %t got %t", want, got)
	}
}

func TestPathSeenTriggerPathImplicitAndOptionalIncludePipeline(t *testing.T) {

	data := `
kind: pipeline
type: docker
name: default

trigger:
  paths:
    - README.md

steps:
- name: message
  image: busybox
  commands:
  - echo "README.md was changed"
`
	pathSeen, err := pathSeen(data)
	if err != nil {
		t.Error(err)
		return
	}

	if !pathSeen {
		t.Errorf("Want 'true' got %t", pathSeen)
	}
}

func TestPathSeenTriggerPathExcludePipeline(t *testing.T) {

	data := `
kind: pipeline
type: docker
name: default

trigger:
  paths:
    exclude:
    - README.md

steps:
- name: message
  image: busybox
  commands:
  - echo "This pipeline will be excluded when README.md is changed"
`

	pathSeen, err := pathSeen(data)
	if err != nil {
		t.Error(err)
		return
	}

	if want, got := true, pathSeen; want != got {
		t.Errorf("Want %t got %t", want, got)
	}
}

func TestPathSeenTriggerPathIncludeAndExcludePipeline(t *testing.T) {

	data := `
kind: pipeline
type: docker
name: default

trigger:
  paths:
    include:
    - CHANGELOG.md
    exclude:
    - README.md

steps:
- name: message
  image: busybox
  commands:
  - echo "CHANGELOG.md was changed, README.md was not changed"
`

	pathSeen, err := pathSeen(data)
	if err != nil {
		t.Error(err)
		return
	}

	if want, got := true, pathSeen; want != got {
		t.Errorf("Want %t got %t", want, got)
	}
}

// parsePipelines tests
func TestParsePipelinesEmptyPipeline(t *testing.T) {
	req := &converter.Request{
		Build: drone.Build{},
		Repo: drone.Repo{
			Slug:   "somewhere/over-the-rainbow",
			Config: ".drone.yml",
		},
	}

	before := `
kind: pipeline
type: docker
name: default
`

	// parsed pipelines don't have a leading newline...
	after := `kind: pipeline
type: docker
name: default
`

	changedFiles := []string{"README.md"}
	resources, err := parsePipelines(before, req.Build, req.Repo, changedFiles)
	if err != nil {
		t.Error(err)
		return
	}

	c, err := marshal(resources)
	if err != nil {
		t.Error(err)
		return
	}
	config := string(c)

	if want, got := after, config; want != got {
		t.Errorf("Want %v got %v", want, got)
	}
}

func TestParsePipelinesStepPathExcludePipeline(t *testing.T) {
	req := &converter.Request{
		Build: drone.Build{},
		Repo: drone.Repo{
			Slug:   "somewhere/over-the-rainbow",
			Config: ".drone.yml",
		},
	}

	before := `
kind: pipeline
type: docker
name: default

steps:
- name: message
  image: busybox
  commands:
  - echo "This step will be excluded when README.md is changed"
  when:
    paths:
      exclude:
      - README.md
`
	// parsed pipelines don't have a leading newline...
	after := `kind: pipeline
type: docker
steps:
- when:
    paths:
      exclude:
      - README.md
    event:
      exclude:
      - '*'
  commands:
  - echo "This step will be excluded when README.md is changed"
  image: busybox
  name: message
name: default
`

	changedFiles := []string{"README.md"}
	resources, err := parsePipelines(before, req.Build, req.Repo, changedFiles)
	if err != nil {
		t.Error(err)
		return
	}

	c, err := marshal(resources)
	if err != nil {
		t.Error(err)
		return
	}
	config := string(c)

	if want, got := after, config; want != got {
		t.Errorf("Want %v got %v", want, got)
	}
}

func TestParsePipelinesStepPathIncludePipeline(t *testing.T) {
	req := &converter.Request{
		Build: drone.Build{},
		Repo: drone.Repo{
			Slug:   "somewhere/over-the-rainbow",
			Config: ".drone.yml",
		},
	}

	before := `
kind: pipeline
type: docker
name: default

steps:
- name: message
  image: busybox
  commands:
  - echo "This step will be included when README.md is changed"
  when:
    paths:
      include:
      - README.md
`

	// parsed pipelines don't have a leading newline...
	after := `kind: pipeline
type: docker
steps:
- when:
    paths:
      include:
      - README.md
  commands:
  - echo "This step will be included when README.md is changed"
  image: busybox
  name: message
name: default
`

	changedFiles := []string{"README.md"}
	resources, err := parsePipelines(before, req.Build, req.Repo, changedFiles)
	if err != nil {
		t.Error(err)
		return
	}

	c, err := marshal(resources)
	if err != nil {
		t.Error(err)
		return
	}
	config := string(c)

	if want, got := after, config; want != got {
		t.Errorf("Want %v got %v", want, got)
	}
}

func TestParsePipelinesStepPathImplicitAndOptionalIncludePipeline(t *testing.T) {
	req := &converter.Request{
		Build: drone.Build{},
		Repo: drone.Repo{
			Slug:   "somewhere/over-the-rainbow",
			Config: ".drone.yml",
		},
	}

	before := `
kind: pipeline
type: docker
name: default

steps:
- name: message
  image: busybox
  commands:
  - echo "README.md was changed"
  when:
    paths:
      - README.md
`

	// parsed pipelines don't have a leading newline...
	after := `kind: pipeline
type: docker
steps:
- when:
    paths:
      include:
      - README.md
  commands:
  - echo "README.md was changed"
  image: busybox
  name: message
name: default
`

	changedFiles := []string{"README.md"}
	resources, err := parsePipelines(before, req.Build, req.Repo, changedFiles)
	if err != nil {
		t.Error(err)
		return
	}

	c, err := marshal(resources)
	if err != nil {
		t.Error(err)
		return
	}
	config := string(c)

	if want, got := after, config; want != got {
		t.Errorf("Want %v got %v", want, got)
	}
}

func TestParsePipelinesTriggerPathImplicitAndOptionalIncludePipeline(t *testing.T) {
	req := &converter.Request{
		Build: drone.Build{},
		Repo: drone.Repo{
			Slug:   "somewhere/over-the-rainbow",
			Config: ".drone.yml",
		},
	}

	before := `
kind: pipeline
type: docker
name: default

trigger:
  paths:
    - README.md

steps:
- name: message
  image: busybox
  commands:
  - echo "README.md was changed"
`

	// parsed pipelines don't have a leading newline...
	after := `kind: pipeline
type: docker
steps:
- commands:
  - echo "README.md was changed"
  image: busybox
  name: message
trigger:
  paths:
    include:
    - README.md
name: default
`

	changedFiles := []string{"README.md"}
	resources, err := parsePipelines(before, req.Build, req.Repo, changedFiles)
	if err != nil {
		t.Error(err)
		return
	}

	c, err := marshal(resources)
	if err != nil {
		t.Error(err)
		return
	}
	config := string(c)

	if want, got := after, config; want != got {
		t.Errorf("Want %v got %v", want, got)
	}
}

func TestParsePipelinesStepPathExcludeAnchorPipeline(t *testing.T) {
	req := &converter.Request{
		Build: drone.Build{},
		Repo: drone.Repo{
			Slug:   "somewhere/over-the-rainbow",
			Config: ".drone.yml",
		},
	}

	before := `
kind: pipeline
type: docker
name: default

commands: &commands
  commands:
  - echo "This step will be excluded when README.md is changed"

steps:
- <<: *commands
  name: message
  image: busybox
  when:
    paths:
      exclude:
      - README.md
`

	after := `kind: pipeline
type: docker
steps:
- when:
    paths:
      exclude:
      - README.md
    event:
      exclude:
      - '*'
  commands:
  - echo "This step will be excluded when README.md is changed"
  image: busybox
  name: message
commands:
  commands:
  - echo "This step will be excluded when README.md is changed"
name: default
`

	changedFiles := []string{"README.md"}
	resources, err := parsePipelines(before, req.Build, req.Repo, changedFiles)
	if err != nil {
		t.Error(err)
		return
	}

	c, err := marshal(resources)
	if err != nil {
		t.Error(err)
		return
	}
	config := string(c)

	if want, got := after, config; want != got {
		t.Errorf("Want %v got %v", want, got)
	}
}

func TestParsePipelinesTriggerPathExcludePipeline(t *testing.T) {
	req := &converter.Request{
		Build: drone.Build{},
		Repo: drone.Repo{
			Slug:   "somewhere/over-the-rainbow",
			Config: ".drone.yml",
		},
	}

	before := `
kind: pipeline
type: docker
name: default

trigger:
  paths:
    exclude:
    - README.md

steps:
- name: message
  image: busybox
  commands:
  - echo "This pipeline will be excluded when README.md is changed"
`

	after := `kind: pipeline
type: docker
steps:
- commands:
  - echo "This pipeline will be excluded when README.md is changed"
  image: busybox
  name: message
trigger:
  paths:
    exclude:
    - README.md
  event:
    exclude:
    - '*'
name: default
`

	changedFiles := []string{"README.md"}
	resources, err := parsePipelines(before, req.Build, req.Repo, changedFiles)
	if err != nil {
		t.Error(err)
		return
	}

	c, err := marshal(resources)
	if err != nil {
		t.Error(err)
		return
	}
	config := string(c)

	if want, got := after, config; want != got {
		t.Errorf("Want %v got %v", want, got)
	}
}

func TestParsePipelinesTriggerPathIncludePipeline(t *testing.T) {
	req := &converter.Request{
		Build: drone.Build{},
		Repo: drone.Repo{
			Slug:   "somewhere/over-the-rainbow",
			Config: ".drone.yml",
		},
	}

	before := `
kind: pipeline
type: docker
name: default

trigger:
  paths:
    include:
    - README.md

steps:
- name: message
  image: busybox
  commands:
  - echo "README.md was changed"
`

	after := `kind: pipeline
type: docker
steps:
- commands:
  - echo "README.md was changed"
  image: busybox
  name: message
trigger:
  paths:
    include:
    - README.md
name: default
`

	changedFiles := []string{"README.md"}
	resources, err := parsePipelines(before, req.Build, req.Repo, changedFiles)
	if err != nil {
		t.Error(err)
		return
	}

	c, err := marshal(resources)
	if err != nil {
		t.Error(err)
		return
	}
	config := string(c)

	if want, got := after, config; want != got {
		t.Errorf("Want %v got %v", want, got)
	}
}

func TestParsePipelinesMultipleTriggerPathIncludePipelines(t *testing.T) {
	req := &converter.Request{
		Build: drone.Build{},
		Repo: drone.Repo{
			Slug:   "somewhere/over-the-rainbow",
			Config: ".drone.yml",
		},
	}

	before := `
kind: pipeline
name: networkredux/staging

node:
  datacenter: nr-staging

concurrency:
  limit: 1

steps:

- name: inspect
  image: meltwaterfoundation/drone-git
  commands:
  - _scripts/inspect.sh networkredux/staging

- name: verify
  image: meltwaterfoundation/drone-lighter:0.21.0
  commands:
  - _scripts/verify.sh
  when:
    event:
    - pull_request
    - push

- name: deploy
  image: meltwaterfoundation/drone-lighter:0.21.0
  commands:
  - _scripts/deploy.sh networkredux/staging
  when:
    event:
    - push

trigger:
  paths:
    include:
    - networkredux/staging/**
  branch:
  - master

---
kind: pipeline
name: basefarm/production

node:
  datacenter: bf-production

concurrency:
  limit: 1

steps:

- name: inspect
  image: meltwaterfoundation/drone-git
  commands:
  - _scripts/inspect.sh basefarm/production

- name: verify
  image: meltwaterfoundation/drone-lighter:0.21.0
  commands:
  - _scripts/verify.sh
  when:
    event:
    - pull_request
    - push

- name: deploy
  image: meltwaterfoundation/drone-lighter:0.21.0
  commands:
  - _scripts/deploy.sh basefarm/production
  when:
    event:
    - push

trigger:
  paths:
    include:
    - basefarm/production/**
  branch:
  - master

...
`

	after := `kind: pipeline
type: ""
steps:
- commands:
  - _scripts/inspect.sh networkredux/staging
  image: meltwaterfoundation/drone-git
  name: inspect
- when:
    event:
    - pull_request
    - push
  commands:
  - _scripts/verify.sh
  image: meltwaterfoundation/drone-lighter:0.21.0
  name: verify
- when:
    event:
    - push
  commands:
  - _scripts/deploy.sh networkredux/staging
  image: meltwaterfoundation/drone-lighter:0.21.0
  name: deploy
trigger:
  paths:
    include:
    - networkredux/staging/**
  branch:
  - master
  event:
    exclude:
    - '*'
concurrency:
  limit: 1
name: networkredux/staging
node:
  datacenter: nr-staging
---
kind: pipeline
type: ""
steps:
- commands:
  - _scripts/inspect.sh basefarm/production
  image: meltwaterfoundation/drone-git
  name: inspect
- when:
    event:
    - pull_request
    - push
  commands:
  - _scripts/verify.sh
  image: meltwaterfoundation/drone-lighter:0.21.0
  name: verify
- when:
    event:
    - push
  commands:
  - _scripts/deploy.sh basefarm/production
  image: meltwaterfoundation/drone-lighter:0.21.0
  name: deploy
trigger:
  paths:
    include:
    - basefarm/production/**
  branch:
  - master
concurrency:
  limit: 1
name: basefarm/production
node:
  datacenter: bf-production
`

	changedFiles := []string{"basefarm/production/foundation/globals.yml"}
	resources, err := parsePipelines(before, req.Build, req.Repo, changedFiles)
	if err != nil {
		t.Error(err)
		return
	}

	c, err := marshal(resources)
	if err != nil {
		t.Error(err)
		return
	}
	config := string(c)

	if want, got := after, config; want != got {
		t.Errorf("Want %v got %v", want, got)
	}
}
