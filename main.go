package main

import (
	"net/http"
	"os"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
	"github.com/zolamk/hasura-exporter/collector"
	"github.com/zolamk/hasura-exporter/settings"
)

func main() {

	logrus.SetFormatter(&logrus.JSONFormatter{})

	logrus.SetOutput(os.Stdout)

	level, _ := logrus.ParseLevel(settings.LogLevel)

	logrus.SetLevel(level)

	registry := prometheus.NewRegistry()

	errors := prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "hasura_errors_total",
		Help: "the total number of errors per collector",
	}, []string{"collector"})

	timeout := time.Duration(10) * time.Second

	registry.MustRegister(errors)

	registry.MustRegister(collector.NewMetadataCollector(errors, timeout))

	registry.MustRegister(collector.NewEventTriggerCollector(errors, timeout))

	registry.MustRegister(collector.NewCronTriggerCollector(errors, timeout))

	http.Handle(settings.WebPath, promhttp.HandlerFor(registry, promhttp.HandlerOpts{}))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`<html>
			<head><title>Hasura Exporter</title></head>
			<body>
			<h1>Hasura Exporter</h1>
			<p><a href="` + settings.WebAddr + `">Metrics</a></p>
			</body>
			</html>`))
	})

	logrus.Info("started hasura exporter on ", settings.WebAddr)

	logrus.WithField("msg", http.ListenAndServe(settings.WebAddr, nil)).Fatalln()

}
