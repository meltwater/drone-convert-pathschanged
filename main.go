// Copyright 2019 the Drone Authors. All rights reserved.
// Use of this source code is governed by the Blue Oak Model License
// that can be found in the LICENSE file.

package main

import (
	"io"
	"net/http"
	"github.com/drone/drone-go/plugin/converter"
	"github.com/meltwater/drone-convert-pathschanged/plugin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	_ "github.com/joho/godotenv/autoload"
	"github.com/kelseyhightower/envconfig"
	"github.com/sirupsen/logrus"

)

// spec provides the plugin settings.
type spec struct {
	Bind   string `envconfig:"DRONE_BIND"`
	Debug  bool   `envconfig:"DRONE_DEBUG"`
	Text   bool   `envconfig:"DRONE_LOGS_TEXT"`
	Secret string `envconfig:"DRONE_SECRET"`

	Token string `envconfig:"GITHUB_TOKEN"`
}

func main() {
	spec := new(spec)
	err := envconfig.Process("", spec)
	if err != nil {
		logrus.Fatal(err)
	}

	if spec.Debug {
		logrus.SetLevel(logrus.DebugLevel)
	}
	if spec.Text {
		logrus.SetFormatter(&logrus.TextFormatter{})
	} else {
		logrus.SetFormatter(&logrus.JSONFormatter{})
	}
	if spec.Secret == "" {
		logrus.Fatalln("missing secret key")
	}
	if spec.Token == "" {
		logrus.Fatalln("missing token")
	}
	if spec.Bind == "" {
		spec.Bind = ":3000"
	}

	handler := converter.Handler(
		plugin.New(
			spec.Token,
		),
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