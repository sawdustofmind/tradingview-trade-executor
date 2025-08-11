package metrics

import "github.com/prometheus/client_golang/prometheus"

const (
	subsystem = "fren"
)

var (
	DurationHistogram = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Subsystem: subsystem,
		Name:      "request_duration_histogram",
		Help:      "Request duration histogram by seconds",
		Buckets:   []float64{.005, .01, .025, .05, .1, .25, .5, 1, 2.5, 5, 10, 15, 20, 30, 40, 50, 60, 120, 180, 600},
	}, []string{"method", "path", "response_status"})

	CustomerCount = prometheus.NewGauge(prometheus.GaugeOpts{
		Subsystem: subsystem,
		Name:      "customer_count",
		Help:      "Total count of customers",
	})

	StrategiesCount = prometheus.NewGauge(prometheus.GaugeOpts{
		Subsystem: subsystem,
		Name:      "strategies_count",
		Help:      "Total count of strategies",
	})

	SubscriptionsCount = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Subsystem: subsystem,
		Name:      "subscriptions_count",
		Help:      "Total count of subscriptions",
	}, []string{"name", "sub_type"})

	SubscriptionsAmount = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Subsystem: subsystem,
		Name:      "subscriptions_amount",
		Help:      "Total amount of subscriptions",
	}, []string{"name", "sub_type"})
)

func init() {
	prometheus.MustRegister(
		DurationHistogram,
		CustomerCount,
		StrategiesCount,
		SubscriptionsCount,
		SubscriptionsAmount,
	)
}
