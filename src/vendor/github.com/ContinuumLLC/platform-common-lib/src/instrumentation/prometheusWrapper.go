package instrumentation

import (
	"fmt"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const (
	port        = ":2112"
	route       = "/metrics"
	libraryName = "GoInstrumentation:"
)

type logger interface {
	Infof(format string, v ...interface{})
	Errorf(format string, v ...interface{})
}

func StartListening(enableListener bool, logger logger) error {
	if logger == nil {
		return fmt.Errorf("The logger passed in nil")
	}
	logger.Infof("%v: value passed %v ", libraryName, enableListener)

	if enableListener {
		logger.Infof("%v: Starting server; go metrics hosted at port: %v, route: %v", libraryName, port, route)
		http.Handle(route, promhttp.Handler())
		err := http.ListenAndServe(port, nil)
		if err != nil {
			return nil
		}
		logger.Infof("%s: Stopping metrics server ", libraryName)
	}
	return nil
}

//RegisterMetric registers a new metric, returns a prometheus.Counter, that you can only increment
func RegisterNewMetricsCounter(name string, help string) prometheus.Counter {
	return promauto.NewCounter(prometheus.CounterOpts{
		Name: name,
		Help: help,
	})
}

//RegisterNewGaugeCounter creates a new gauge metric, returns Prometheus.Guage, that you can increment or decrement in the consuming code
func RegisterNewGaugeCounter(name string, help string) prometheus.Gauge {
	return promauto.NewGauge(prometheus.GaugeOpts{
		Name: name,
		Help: help,
	})
}

func NewBucket(start, width float64, count int) Bucket {
	return Bucket{
		Start: start,
		Width: width,
		Count: count,
	}
}

type Bucket struct {
	Start float64
	Width float64
	Count int
}

//RegisterNewHistrogram Summary captures individual observations from an event or sample stream and summarizes them in a manner similar to traditional summary statistics: 1. sum of observations, 2. observation count, 3. rank estimations.
func RegisterNewHistrogram(name string, help string, bucket Bucket) prometheus.Histogram {
	return promauto.NewHistogram(prometheus.HistogramOpts{
		Name:    name,
		Help:    help,
		Buckets: prometheus.LinearBuckets(bucket.Start, bucket.Width, bucket.Count),
	})
}

//RegisterNewSummary created a new Summary: A Summary captures individual observations from an event or sample stream and summarizes them in a manner similar to traditional summary statistics: 1. sum of observations, 2. observation count, 3. rank estimations.
func RegisterNewSummary(name string, help string, objectives map[float64]float64) prometheus.Summary {
	return promauto.NewSummary(prometheus.SummaryOpts{
		Name:       name,
		Help:       help,
		Objectives: objectives,
	})
}
