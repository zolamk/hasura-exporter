package collector

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/sirupsen/logrus"
	"github.com/zolamk/hasura-exporter/settings"
)

type MetadataCollector struct {
	errors       *prometheus.CounterVec
	timeout      time.Duration
	Inconsistent *prometheus.Desc
}

func NewMetadataCollector(errors *prometheus.CounterVec, timeout time.Duration) *MetadataCollector {

	labels := []string{}

	return &MetadataCollector{
		errors:  errors,
		timeout: timeout,
		Inconsistent: prometheus.NewDesc(
			"hasura_metadata_consistency_status",
			"If 1 metadata is consistent, 0 otherwise",
			labels,
			nil,
		),
	}
}

func (c *MetadataCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.Inconsistent
}

func (c *MetadataCollector) Collect(ch chan<- prometheus.Metric) {

	var err error

	metadata_api_url := fmt.Sprintf("%s/v1/metadata", settings.HasuraGraphQLURL)

	body, _ := json.Marshal(map[string]interface{}{
		"type": "get_inconsistent_metadata",
		"args": map[string]interface{}{},
	})

	req, err := http.NewRequest("POST", metadata_api_url, bytes.NewReader(body))

	if err != nil {
		c.errors.WithLabelValues("metadata").Add(1)
		logrus.WithField("err", err).Error("can't get metadata status")
		return
	}

	req.Header.Add("content-type", "application/json")

	req.Header.Add("x-hasura-admin-secret", settings.HasuraAdminSecret)

	res, err := http.DefaultClient.Do(req)

	if err != nil {
		c.errors.WithLabelValues("metadata").Add(1)
		logrus.WithField("err", err).Error("can't get metadata status")
		return
	}

	if res.StatusCode != http.StatusOK {

		var body []byte

		res.Body.Read(body)

		c.errors.WithLabelValues("metadata").Add(1)

		logrus.WithField("status_code", res.StatusCode).WithField("response_body", string(body)).Error("can't get metadata status")

		return

	}

	var data struct {
		IsConsistent        bool                     `json:"is_consistent"`
		InconsistentObjects []map[string]interface{} `json:"inconsistent_objects"`
	}

	decoder := json.NewDecoder(res.Body)

	err = decoder.Decode(&data)

	if err != nil {
		c.errors.WithLabelValues("metadata").Add(1)
		logrus.WithField("err", err).Error("can't get metadata status")
		return
	}

	metadata_status := 0.0

	if data.IsConsistent {
		metadata_status = 1.0
	}

	ch <- prometheus.MustNewConstMetric(
		c.Inconsistent,
		prometheus.GaugeValue,
		metadata_status,
	)

}
