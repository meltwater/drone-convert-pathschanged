package plugin

import (
	"context"
	"testing"

	"github.com/drone/drone-go/drone"
	"github.com/drone/drone-go/plugin/converter"
)

// empty context
var noContext = context.Background()

func TestPluginEmptyPipeline(t *testing.T) {

	providers := []string{"github","bitbucket-server"}

	req := &converter.Request{
		Build: drone.Build{},
		Repo: drone.Repo{
			Slug:   "somewhere/over-the-rainbow",
			Config: ".drone.yml",
		},
	}

	for _, provider := range providers {
		plugin := New("invalidtoken", provider)

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
