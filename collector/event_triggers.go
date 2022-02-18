package collector

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/sirupsen/logrus"
	"github.com/zolamk/hasura-exporter/settings"
)

type EventTriggerCollector struct {
	errors     *prometheus.CounterVec
	timeout    time.Duration
	Pending    *prometheus.Desc
	Processed  *prometheus.Desc
	Successful *prometheus.Desc
	Failed     *prometheus.Desc
}

func NewEventTriggerCollector(errors *prometheus.CounterVec, timeout time.Duration) *EventTriggerCollector {

	labels := []string{
		"trigger_name",
	}

	return &EventTriggerCollector{
		errors:  errors,
		timeout: timeout,
		Pending: prometheus.NewDesc(
			"hasura_pending_event_triggers",
			"number of pending hasura event triggers",
			labels,
			nil,
		),
		Processed: prometheus.NewDesc(
			"hasura_processed_event_triggers",
			"number of processed hasura event triggers",
			labels,
			nil,
		),
		Successful: prometheus.NewDesc(
			"hasura_successful_event_triggers",
			"number of successfully processed hasura event triggers",
			labels,
			nil,
		),
		Failed: prometheus.NewDesc(
			"hasura_failed_event_triggers",
			"number of failed hasura event triggers",
			labels,
			nil,
		),
	}
}

func (c *EventTriggerCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.Pending
	ch <- c.Processed
	ch <- c.Successful
	ch <- c.Failed
}

func (c *EventTriggerCollector) Collect(ch chan<- prometheus.Metric) {

	var err error

	query_api_url := fmt.Sprintf("%s/v2/query", settings.HasuraGraphQLEndpoint)

	body, _ := json.Marshal(map[string]interface{}{
		"type": "bulk",
		"args": []map[string]interface{}{
			{
				"type": "run_sql",
				"args": map[string]interface{}{
					"cascade":   false,
					"read_only": true,
					"sql":       "SELECT COUNT(*), trigger_name FROM hdb_catalog.event_log WHERE delivered = true OR error = true GROUP BY trigger_name;",
				},
			},
			{
				"type": "run_sql",
				"args": map[string]interface{}{
					"cascade":   false,
					"read_only": true,
					"sql":       "SELECT COUNT(*), trigger_name FROM hdb_catalog.event_log WHERE delivered = false AND error = false AND archived = false GROUP BY trigger_name;",
				},
			},
			{
				"type": "run_sql",
				"args": map[string]interface{}{
					"cascade":   false,
					"read_only": true,
					"sql":       "SELECT COUNT(*), trigger_name FROM hdb_catalog.event_log WHERE error = true GROUP BY trigger_name;",
				},
			},
			{
				"type": "run_sql",
				"args": map[string]interface{}{
					"cascade":   false,
					"read_only": true,
					"sql":       "SELECT COUNT(*), trigger_name FROM hdb_catalog.event_log WHERE error = false AND delivered = true GROUP BY trigger_name;",
				},
			},
		},
	})

	req, err := http.NewRequest("POST", query_api_url, bytes.NewReader(body))

	if err != nil {

		c.errors.WithLabelValues("event").Add(1)

		logrus.WithField("err", err).Error("can't get event data")

		return

	}

	req.Header.Add("content-type", "application/json")

	req.Header.Add("x-hasura-admin-secret", settings.HasuraAdminSecret)

	res, err := http.DefaultClient.Do(req)

	if err != nil {

		c.errors.WithLabelValues("event").Add(1)

		logrus.WithField("err", err).Error("can't get event data")

		return

	}

	if res.StatusCode != http.StatusOK {

		var body []byte

		res.Body.Read(body)

		c.errors.WithLabelValues("event").Add(1)

		logrus.WithField("status_code", res.StatusCode).WithField("response_body", string(body)).Error("can't get event data")

		return

	}

	var data []struct {
		ResultType string     `json:"result_type"`
		Result     [][]string `json:"result"`
	}

	decoder := json.NewDecoder(res.Body)

	err = decoder.Decode(&data)

	if err != nil {

		c.errors.WithLabelValues("event").Add(1)

		logrus.WithField("err", err).Error("can't get event data")

		return

	}

	for i, result := range data[0].Result {

		if i == 0 {
			continue
		}

		v, err := strconv.ParseFloat(result[0], 64)

		if err != nil {

			c.errors.WithLabelValues("event").Add(1)

			logrus.WithField("err", err).Error("error converting count string to float")

			continue

		}

		ch <- prometheus.MustNewConstMetric(
			c.Processed,
			prometheus.GaugeValue,
			v,
			result[1],
		)

	}

	for i, result := range data[1].Result {

		if i == 0 {
			continue
		}

		v, err := strconv.ParseFloat(result[0], 64)

		if err != nil {

			c.errors.WithLabelValues("event").Add(1)

			logrus.WithField("err", err).Error("error converting count string to float")

			continue

		}

		ch <- prometheus.MustNewConstMetric(
			c.Pending,
			prometheus.GaugeValue,
			v,
			result[1],
		)

	}

	for i, result := range data[2].Result {

		if i == 0 {
			continue
		}

		v, err := strconv.ParseFloat(result[0], 64)

		if err != nil {

			c.errors.WithLabelValues("event").Add(1)

			logrus.WithField("err", err).Error("error converting count string to float")

			continue

		}

		ch <- prometheus.MustNewConstMetric(
			c.Failed,
			prometheus.GaugeValue,
			v,
			result[1],
		)

	}

	for i, result := range data[3].Result {

		if i == 0 {
			continue
		}

		v, err := strconv.ParseFloat(result[0], 64)

		if err != nil {

			c.errors.WithLabelValues("event").Add(1)

			logrus.WithField("err", err).Error("error converting count string to float")

			continue

		}

		ch <- prometheus.MustNewConstMetric(
			c.Successful,
			prometheus.GaugeValue,
			v,
			result[1],
		)

	}

}
