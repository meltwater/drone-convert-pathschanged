package plugin

import (
	"context"
	"io/ioutil"
	"testing"

	"github.com/drone/drone-go/drone"
	"github.com/drone/drone-go/plugin/converter"
)

// empty context
var noContext = context.Background()

func TestPluginEmptyPipeline(t *testing.T) {
	req := &converter.Request{
		Build: drone.Build{},
		Repo: drone.Repo{
			Slug:   "somewhere/over-the-rainbow",
			Config: ".drone.yml",
		},
	}

	plugin := New("invalidtoken")

	config, err := plugin.Convert(noContext, req)
	if err != nil {
		t.Error(err)
		return
	}

	if want, got := config.Data, ""; want != got {
		t.Errorf("Want %q got %q", want, got)
	}
}

func TestParsePipelinesEmptyPipeline(t *testing.T) {
	req := &converter.Request{
		Build: drone.Build{},
		Repo: drone.Repo{
			Slug:   "somewhere/over-the-rainbow",
			Config: ".drone.yml",
		},
	}

	before, err := ioutil.ReadFile("testdata/single_empty_pipeline.yml")
	if err != nil {
		t.Error(err)
		return
	}

	after, err := ioutil.ReadFile("testdata/single_empty_pipeline.yml.golden")
	if err != nil {
		t.Error(err)
		return
	}

	data := string(before)

	changedFiles := []string{"README.md"}
	resources, err := parsePipelines(data, req.Build, req.Repo, changedFiles)
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

	if want, got := string(after), config; want != got {
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

	before, err := ioutil.ReadFile("testdata/single_step_with_exclude_pipeline.yml")
	if err != nil {
		t.Error(err)
		return
	}

	after, err := ioutil.ReadFile("testdata/single_step_with_exclude_pipeline.yml.golden")
	if err != nil {
		t.Error(err)
		return
	}

	data := string(before)

	changedFiles := []string{"README.md"}
	resources, err := parsePipelines(data, req.Build, req.Repo, changedFiles)
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

	if want, got := string(after), config; want != got {
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

	before, err := ioutil.ReadFile("testdata/single_step_with_exclude_anchor_pipeline.yml")
	if err != nil {
		t.Error(err)
		return
	}

	after, err := ioutil.ReadFile("testdata/single_step_with_exclude_anchor_pipeline.yml.golden")
	if err != nil {
		t.Error(err)
		return
	}

	data := string(before)

	changedFiles := []string{"README.md"}
	resources, err := parsePipelines(data, req.Build, req.Repo, changedFiles)
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

	if want, got := string(after), config; want != got {
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

	before, err := ioutil.ReadFile("testdata/single_trigger_with_exclude_pipeline.yml")
	if err != nil {
		t.Error(err)
		return
	}

	after, err := ioutil.ReadFile("testdata/single_trigger_with_exclude_pipeline.yml.golden")
	if err != nil {
		t.Error(err)
		return
	}

	data := string(before)

	changedFiles := []string{"README.md"}
	resources, err := parsePipelines(data, req.Build, req.Repo, changedFiles)
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

	if want, got := string(after), config; want != got {
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

	before, err := ioutil.ReadFile("testdata/single_trigger_with_include_pipeline.yml")
	if err != nil {
		t.Error(err)
		return
	}

	after, err := ioutil.ReadFile("testdata/single_trigger_with_include_pipeline.yml.golden")
	if err != nil {
		t.Error(err)
		return
	}

	data := string(before)

	changedFiles := []string{"README.md"}
	resources, err := parsePipelines(data, req.Build, req.Repo, changedFiles)
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

	if want, got := string(after), config; want != got {
		t.Errorf("Want %v got %v", want, got)
	}
}
