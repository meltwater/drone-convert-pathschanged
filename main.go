// Copyright 2019 the Drone Authors. All rights reserved.
// Use of this source code is governed by the Blue Oak Model License
// that can be found in the LICENSE file.

package main

import (
	"flag"
	"fmt"
	"io"
	"net/http"

	"github.com/meltwater/drone-convert-pathschanged/plugin"

	"github.com/drone/drone-go/plugin/converter"
	"github.com/joho/godotenv"
	_ "github.com/joho/godotenv/autoload"
	"github.com/kelseyhightower/envconfig"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
)

// spec provides the plugin settings.
type (
	spec struct {
		Bind   string `envconfig:"DRONE_BIND"`
		Debug  bool   `envconfig:"DRONE_DEBUG"`
		Text   bool   `envconfig:"DRONE_LOGS_TEXT"`
		Secret string `envconfig:"DRONE_SECRET"`

		Provider string `envconfig:"PROVIDER"`
		Token    string `envconfig:"TOKEN"`
		// BB_ADDRESS is deprecated in favor of STASH_SERVER, it will be removed in a future version
		BitBucketAddress  string `envconfig:"BB_ADDRESS"`
		BitBucketUser     string `envconfig:"BITBUCKET_USER"`
		BitBucketPassword string `envconfig:"BITBUCKET_PASSWORD"`
		GithubServer      string `envconfig:"GITHUB_SERVER"`
		StashServer       string `envconfig:"STASH_SERVER"`
	}
)

func contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}
	return false
}

func validate(spec *spec) error {
	if spec.Secret == "" {
		return fmt.Errorf("missing secret key")
	}
	if spec.Provider == "" {
		return fmt.Errorf("missing provider")
	} else {
		providers := []string{
			"bitbucket",
			// bitbucket-server support is deprecated in favor of stash, it will be removed in a future version
			"bitbucket-server",
			"github",
			"stash",
			"gitee",
		}
		if !contains(providers, spec.Provider) {
			return fmt.Errorf("unsupported provider")
		}
	}
	if spec.Token == "" && (spec.Provider == "github" || spec.Provider == "bitbucket-server" || spec.Provider == "stash") {
		return fmt.Errorf("missing token")
	}
	if spec.BitBucketUser == "" && spec.Provider == "bitbucket" {
		return fmt.Errorf("missing bitbucket user")
	}
	if spec.BitBucketPassword == "" && spec.Provider == "bitbucket" {
		return fmt.Errorf("missing bitbucket password")
	}
	if spec.BitBucketAddress == "" && spec.Provider == "bitbucket-server" {
		return fmt.Errorf("missing bitbucket server address")
	} else if spec.BitBucketAddress != "" && spec.Provider == "bitbucket-server" {
		// backwards compatible support for bitbucket-server, this will be removed in a future version
		spec.StashServer = spec.BitBucketAddress
		spec.Provider = "stash"

		logrus.Warningln("bitbucket-server support is deprecated, please use stash")
	}
	if spec.StashServer == "" && spec.Provider == "stash" {
		return fmt.Errorf("missing stash server")
	}

	if spec.Token == "" && spec.Provider == "gitee" {
		return fmt.Errorf("missing gitee token")
	}

	return nil
}

func main() {
	envfilepath := flag.String("envfile", "", "pass filepath to env file (optional)")
	flag.Parse()
	if *envfilepath != "" {
		godotenv.Load(*envfilepath)
	}
	spec := new(spec)
	err := envconfig.Process("", spec)
	if err != nil {
		logrus.Fatalln(err)
	}

	if spec.Debug {
		logrus.SetLevel(logrus.DebugLevel)
	}
	if spec.Text {
		logrus.SetFormatter(&logrus.TextFormatter{})
	} else {
		logrus.SetFormatter(&logrus.JSONFormatter{})
	}

	err = validate(spec)
	if err != nil {
		logrus.Fatalln(err)
	}

	if spec.Bind == "" {
		spec.Bind = ":3000"
	}

	params := &plugin.Params{
		BitBucketUser:     spec.BitBucketUser,
		BitBucketPassword: spec.BitBucketPassword,
		GithubServer:      spec.GithubServer,
		Token:             spec.Token,
		StashServer:       spec.StashServer,
	}

	handler := converter.Handler(
		plugin.New(spec.Provider, params),
		spec.Secret,
		logrus.StandardLogger(),
	)

	logrus.Infof("server listening on address %s", spec.Bind)

	http.Handle("/", handler)
	http.HandleFunc("/healthz", healthz)
	http.Handle("/metrics", promhttp.Handler())
	logrus.Fatal(http.ListenAndServe(spec.Bind, nil))
}

func healthz(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	io.WriteString(w, "OK")
}
