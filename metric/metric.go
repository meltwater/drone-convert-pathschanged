package metric

import (
        "github.com/prometheus/client_golang/prometheus"
        "github.com/prometheus/client_golang/prometheus/promauto"
)

var (
        GithubApiCount = promauto.NewCounter(prometheus.CounterOpts{
                Name: "github_api_calls",
                Help: "Total number of github api calls per hour",
		})
)