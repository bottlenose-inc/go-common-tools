package metrics

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/bottlenose-inc/go-common-tools/logger"      // go-common-tools logger package
	"github.com/prometheus/client_golang/prometheus" // Official Prometheus golang library
)

type PrometheusId struct {
	Name    string
	Address string
	Port    int
	ID      string
}

var (
	histogramBuckets = []float64{0.001, 0.0025, 0.005, 0.0075, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1, 2.5, 5, 10, 20, 30, 45, 60, 90}
)

func StartPrometheusMetricsServer(name string, logger *logger.Logger, port int) error {
	// name for identifying the service
	// logger - Logger object from go-common-tools#logger.go
	// port for Prometheus to report metrics to
	// Returns an error or nil upon successful setup

	// Start HTTP server
	http.Handle("/metrics", prometheus.Handler())
	err := http.ListenAndServe(":"+strconv.Itoa(port), nil)
	if err != nil {
		logger.Error("Error starting Prometheus metrics server: " + err.Error())
		return err
	}
	return nil
}

func CreateHistogram(name string, namespace string, subsystem string, help string, labels map[string]string, buckets ...[]float64) (histogram prometheus.Histogram, err error) {
	// "name" and "help" are required by Prometheus to create a histogram
	// all other fields are optional
	// Returns a prometheus histogram object

	useBuckets := histogramBuckets
	if len(buckets) > 0 {
		useBuckets = buckets[0]
	}

	constLabels := prometheus.Labels(labels)
	if name == "" || help == "" {
		err = errors.New("Prometheus histogram requires both name and help fields to initialize - missing one or both of those fields")
		return nil, err
	}
	histogram = prometheus.NewHistogram(prometheus.HistogramOpts{
		Name:        name,
		Help:        help,
		Namespace:   namespace,
		Subsystem:   subsystem,
		ConstLabels: constLabels,
		Buckets:     useBuckets,
	})

	prometheus.MustRegister(histogram)

	return histogram, nil

}

func CreateHistogramVector(name string, namespace string, subsystem string, help string, labels map[string]string, labelNames []string, buckets ...[]float64) (histogramVec *prometheus.HistogramVec, err error) {
	// "name" and "help" are required by Prometheus to create a histogram
	// all other fields are optional
	// Returns a prometheus histogram object

	useBuckets := histogramBuckets
	if len(buckets) > 0 {
		useBuckets = buckets[0]
	}

	constLabels := prometheus.Labels(labels)
	if name == "" || help == "" {
		err = errors.New("Prometheus histogram requires both name and help fields to initialize - missing one or both of those fields")
		return nil, err
	}
	histogramVec = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name:        name,
		Help:        help,
		Namespace:   namespace,
		Subsystem:   subsystem,
		ConstLabels: constLabels,
		Buckets:     useBuckets,
	}, labelNames)

	prometheus.MustRegister(histogramVec)

	return histogramVec, nil

}
func CreateCounterVector(name string, namespace string, subsystem string, help string, labels map[string]string, labelNames []string) (counterVec *prometheus.CounterVec, err error) {
	// "name" and "help" are required by Prometheus to create a counter vector
	// all other fields are optional
	// Returns a prometheus counter vector object

	constLabels := prometheus.Labels(labels)

	if name == "" || help == "" {
		err = errors.New("Prometheus counter vector requires both name and help fields to initialize - missing one or both of those fields")
		return nil, err
	}
	counterVec = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name:        name,
		Help:        help,
		Namespace:   namespace,
		Subsystem:   subsystem,
		ConstLabels: constLabels,
	}, labelNames)

	prometheus.MustRegister(counterVec)

	return counterVec, nil
}

// initialize a counter vector with given labels, setting values to 0
func InitCounterVector(counterVec *prometheus.CounterVec, labels []string) {
	for _, label := range labels {
		counter, err := counterVec.GetMetricWithLabelValues(label)
		if err == nil {
			counter.Add(0)
		}
	}
}

func CreateCounter(name string, namespace string, subsystem string, help string, labels map[string]string) (counter prometheus.Counter, err error) {
	// "name" and "help" are required by Prometheus to create a counter
	// all other fields are optional
	// Returns a prometheus counter object

	constLabels := prometheus.Labels(labels)

	if name == "" || help == "" {
		err = errors.New("Prometheus counter requires both name and help fields to initialize - missing one or both of those fields")
		return nil, err
	}

	counter = prometheus.NewCounter(prometheus.CounterOpts{
		Name:        name,
		Help:        help,
		Namespace:   namespace,
		Subsystem:   subsystem,
		ConstLabels: constLabels,
	})

	prometheus.MustRegister(counter)

	return counter, nil
}

func CreateGauge(name string, namespace string, subsystem string, help string, labels map[string]string) (gauge prometheus.Gauge, err error) {
	// "name" and "help" are required by Prometheus to create a gauge
	// all other fields are optional
	// Returns a prometheus gauge object

	constLabels := prometheus.Labels(labels)

	if name == "" || help == "" {
		err = errors.New("Prometheus gauge requires both name and help fields to initialize - missing one or both of those fields")
		return nil, err
	}

	gauge = prometheus.NewGauge(prometheus.GaugeOpts{
		Name:        name,
		Help:        help,
		Namespace:   namespace,
		Subsystem:   subsystem,
		ConstLabels: constLabels,
	})

	prometheus.MustRegister(gauge)

	return gauge, nil
}

func CreateGaugeVector(name string, namespace string, subsystem string, help string, labels map[string]string, labelNames []string) (gaugeVec *prometheus.GaugeVec, err error) {
	// "name" and "help" are required by Prometheus to create a gauge vector
	// all other fields are optional
	// Returns a prometheus gauge vector object

	constLabels := prometheus.Labels(labels)

	if name == "" || help == "" {
		err = errors.New("Prometheus gauge vector requires both name and help fields to initialize - missing one or both of those fields")
		return nil, err
	}

	gaugeVec = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name:        name,
		Help:        help,
		Namespace:   namespace,
		Subsystem:   subsystem,
		ConstLabels: constLabels,
	}, labelNames)

	prometheus.MustRegister(gaugeVec)

	return gaugeVec, nil
}
