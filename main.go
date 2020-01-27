// Copyright 2019 the Drone Authors. All rights reserved.
// Use of this source code is governed by the Blue Oak Model License
// that can be found in the LICENSE file.

package main

import (
	"io"
	"context"
	"net/http"
    "github.com/google/go-github/github"
	"golang.org/x/oauth2"
	"github.com/drone/drone-go/plugin/converter"
	"github.com/prometheus/client_golang/prometheus"
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

var (
	    githubApiCount = prometheus.NewGauge(
	    prometheus.GaugeOpts{
			   Name: "github_api_calls_remaining",
			   Help: "Total number of github api calls per hour remaining",
	},)
) 
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
	ApiRateLimit(spec.Token)
	
	http.Handle("/", handler)
	http.HandleFunc("/healthz", healthz)
	http.Handle("/metrics", promhttp.Handler())
	prometheus.MustRegister(githubApiCount)
	logrus.Fatal(http.ListenAndServe(spec.Bind, nil))
}

func healthz(w http.ResponseWriter, r *http.Request) {
 	w.Header().Set("Content-Type", "text/plain")
	io.WriteString(w, "OK")
}

func ApiRateLimit(token string) {
	go func() {
		newctx := context.Background()
		ts := oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: token},
		)
		tc := oauth2.NewClient(newctx, ts)
		
		client := github.NewClient(tc)
			
		rateLimit,_,err:= client.RateLimits(newctx)
		if err != nil {
			logrus.Fatalln("No metrics")
		}
		githubApiCount.Set(float64(rateLimit.Core.Remaining))

	}()
}