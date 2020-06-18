// Copyright 2019 the Drone Authors. All rights reserved.
// Use of this source code is governed by the Blue Oak Model License
// that can be found in the LICENSE file.

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

func TestPlugin(t *testing.T) {
	req := &converter.Request{
		Build: drone.Build{
			After: "3d21ec53a331a6f037a91c368710b99387d012c1",
		},
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

	if config.Data != "" {
		t.Error("An empty pipeline returned data")
	}
	before, err := ioutil.ReadFile("testdata/single_pipeline.yml")
	if err != nil {
		t.Error(err)
		return
	}
	after, err := ioutil.ReadFile("testdata/single_pipeline.yml.golden")
	if err != nil {
		t.Error(err)
		return
	}
	req.Repo.Config = "single_pipeline.yml"
	req.Config.Data = string(before)
	config, err = plugin.Convert(noContext, req)
	if err != nil {
		t.Error(err)
		return
	}
	if config.Data == "" {
		t.Error("Want non-empty configuration")
		return
	}
	if want, got := config.Data, string(after); want != got {
		t.Errorf("Want %q got %q", want, got)
	}
}
