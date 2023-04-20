package metrics

import (
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
)

func Start() error {
	http.Handle("/metrics", promhttp.Handler())
	return http.ListenAndServe(":2112", nil)
}
