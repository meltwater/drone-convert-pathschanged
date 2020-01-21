package metric

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/meltwater/drone-convert-pathschanged/plugin"
)

func GithubApiCalls(get plugin.apiRateLimit)  {
	prometheus.MustRegister(
			prometheus.NewGaugeFunc(prometheus.GaugeOpts{
					Name: "github_api_calls"
					Help: "Total number of github api calls per hour."
		   }, func() float64{
				   i, _ := get.apiRateLimit(noContext)
				   return float64(i)
		   }),
		   
	)

}