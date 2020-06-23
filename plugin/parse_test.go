package plugin

import (
	"io/ioutil"
	"testing"

	"github.com/drone/drone-go/drone"
	"github.com/drone/drone-go/plugin/converter"
)

// pathSeen tests
func TestPathSeenEmptyPipeline(t *testing.T) {
	before, err := ioutil.ReadFile("testdata/single_empty_pipeline.yml")
	if err != nil {
		t.Error(err)
		return
	}

	data := string(before)
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
	before, err := ioutil.ReadFile("testdata/single_invalid_pipeline.yml")
	if err != nil {
		t.Error(err)
		return
	}

	data := string(before)
	_, err = pathSeen(data)

	if err == nil {
		t.Errorf("Invalid pipeline did not return an error")
	}
}

func TestPathSeenSingleStepPipeline(t *testing.T) {
	before, err := ioutil.ReadFile("testdata/single_step_pipeline.yml")
	if err != nil {
		t.Error(err)
		return
	}

	data := string(before)
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
	before, err := ioutil.ReadFile("testdata/single_step_with_include_pipeline.yml")
	if err != nil {
		t.Error(err)
		return
	}

	data := string(before)
	pathSeen, err := pathSeen(data)
	if err != nil {
		t.Error(err)
		return
	}

	if want, got := true, pathSeen; want != got {
		t.Errorf("Want %t got %t", want, got)
	}
}

func TestPathSeenStepPathExcludePipeline(t *testing.T) {
	before, err := ioutil.ReadFile("testdata/single_step_with_exclude_pipeline.yml")
	if err != nil {
		t.Error(err)
		return
	}

	data := string(before)
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
	before, err := ioutil.ReadFile("testdata/single_step_with_include_and_exclude_pipeline.yml")
	if err != nil {
		t.Error(err)
		return
	}

	data := string(before)
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
	before, err := ioutil.ReadFile("testdata/single_trigger_with_include_pipeline.yml")
	if err != nil {
		t.Error(err)
		return
	}

	data := string(before)
	pathSeen, err := pathSeen(data)
	if err != nil {
		t.Error(err)
		return
	}

	if want, got := true, pathSeen; want != got {
		t.Errorf("Want %t got %t", want, got)
	}
}

func TestPathSeenTriggerPathExcludePipeline(t *testing.T) {
	before, err := ioutil.ReadFile("testdata/single_trigger_with_exclude_pipeline.yml")
	if err != nil {
		t.Error(err)
		return
	}

	data := string(before)
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
	before, err := ioutil.ReadFile("testdata/single_trigger_with_include_and_exclude_pipeline.yml")
	if err != nil {
		t.Error(err)
		return
	}

	data := string(before)
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

func TestParsePipelinesMultipleTriggerPathIncludePipelines(t *testing.T) {
	req := &converter.Request{
		Build: drone.Build{},
		Repo: drone.Repo{
			Slug:   "somewhere/over-the-rainbow",
			Config: ".drone.yml",
		},
	}

	before, err := ioutil.ReadFile("testdata/multiple_trigger_with_include_pipelines.yml")
	if err != nil {
		t.Error(err)
		return
	}

	after, err := ioutil.ReadFile("testdata/multiple_trigger_with_include_pipelines.yml.golden")
	if err != nil {
		t.Error(err)
		return
	}

	data := string(before)

	changedFiles := []string{"basefarm/production/foundation/globals.yml"}
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
