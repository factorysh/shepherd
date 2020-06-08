package metrics

import (
	"fmt"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
)

var (
	ContainerDead = promauto.NewCounter(prometheus.CounterOpts{
		Name: "container_dead",
		Help: "Amount of dead containers",
	})
	ContainerDestroyed = promauto.NewCounter(prometheus.CounterOpts{
		Name: "container_destroyed",
		Help: "Amount of destroyed containers",
	})
)

func ListenAndServe(listen string) error {
	http.Handle("/metrics", promhttp.Handler())
	log.WithFields(
		log.Fields{
			"url": fmt.Sprintf("http://%s/metrics", listen),
		},
	).Info("metrics.ListenAndServe")
	return http.ListenAndServe(listen, nil)
}
