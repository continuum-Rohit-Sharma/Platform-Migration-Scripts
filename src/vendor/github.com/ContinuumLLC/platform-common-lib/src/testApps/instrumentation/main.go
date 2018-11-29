package main

import (
	"fmt"
	"math"

	"github.com/ContinuumLLC/platform-common-lib/src/instrumentation"
)

func main() {

	go instrumentation.StartListening(true, logger{})
	summaryCounter := instrumentation.RegisterNewSummary("testSummary", "testsummary", map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001})
	histSummaryCounter := instrumentation.RegisterNewHistrogram("testHistogramSummary", "testHistogramSummary", instrumentation.NewBucket(20, 5, 5))
	counter := instrumentation.RegisterNewMetricsCounter("testCounter", "this is a test counter")
	guageCounter := instrumentation.RegisterNewGaugeCounter("testGuageCouter", "this is a test guage counter")

	i := 0
	for {
		guageCounter.Inc()
		counter.Inc()
		histSummaryCounter.Observe(30 + math.Floor(120*math.Sin(float64(i)*0.1))/10)
		summaryCounter.Observe(30 + math.Floor(120*math.Sin(float64(i)*0.1))/10)
		i = 1

	}

}

type logger struct {
}

func (logger) Infof(format string, v ...interface{}) {
	fmt.Printf(format+"\n", v)
}

func (logger) Errorf(format string, v ...interface{}) {
	fmt.Printf(format+"\n", v)
}
