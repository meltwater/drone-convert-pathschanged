package plugin

import (
	"io/ioutil"
	"testing"
)

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
