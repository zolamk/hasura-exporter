package collector

import (
	"fmt"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/sirupsen/logrus"
	"github.com/zolamk/hasura-exporter/settings"
)

type HealthcheckCollector struct {
	errors  *prometheus.CounterVec
	timeout time.Duration
	Healthy *prometheus.Desc
}

func NewHealthcheckCollector(errors *prometheus.CounterVec, timeout time.Duration) *HealthcheckCollector {

	labels := []string{}

	return &HealthcheckCollector{
		errors:  errors,
		timeout: timeout,
		Healthy: prometheus.NewDesc(
			"hasura_healthy",
			"If 1 hasura graphql server is healthy, 0 otherwise",
			labels,
			nil,
		),
	}
}

func (c *HealthcheckCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.Healthy
}

func (c *HealthcheckCollector) Collect(ch chan<- prometheus.Metric) {

	var err error

	health_check_url := fmt.Sprintf("%s/healthz", settings.HasuraGraphQLEndpoint)

	req, err := http.NewRequest("GET", health_check_url, nil)

	if err != nil {
		c.errors.WithLabelValues("health").Add(1)
		logrus.WithField("err", err).Error("can't get health status")
		return
	}

	res, err := http.DefaultClient.Do(req)

	if err != nil {
		c.errors.WithLabelValues("health").Add(1)
		logrus.WithField("err", err).Error("can't get health status")
		return
	}

	health_status := 1.0

	if res.StatusCode == http.StatusInternalServerError {

		health_status = 0.0

	} else if res.StatusCode != http.StatusOK {

		c.errors.WithLabelValues("health").Add(1)

		logrus.WithField("status_code", res.StatusCode).Error("can't get metadata status")

		return

	}

	ch <- prometheus.MustNewConstMetric(
		c.Healthy,
		prometheus.GaugeValue,
		health_status,
	)

}
