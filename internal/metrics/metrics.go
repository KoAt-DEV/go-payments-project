package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

const namespace = "gopayments"
const subsystem = "auth"

var TotalHealthcheck = promauto.NewCounter(
	prometheus.CounterOpts{
		Namespace: namespace,
		Name:      "total_healthchecks",
		Help:      "Total number of healthchecks",
	},
)

var RateLimitExceeded = promauto.NewCounter(
	prometheus.CounterOpts{
		Namespace: namespace,
		Name:      "total_ratelimit_exceeded",
		Help:      "Total requests where rate limit exceeded",
	},
)

var RegisterRequestTotal = promauto.NewCounterVec(
	prometheus.CounterOpts{
		Namespace: namespace,
		Subsystem: subsystem,
		Name:      "total_register_requests",
		Help:      "Total number of register requests",
	},
	[]string{"outcome"},
)

var TotalRegisterRequestTime = promauto.NewHistogram(
	prometheus.HistogramOpts{
		Namespace: namespace,
		Subsystem: subsystem,
		Name:      "total_register_request_time",
		Help:      "Total register request latency",
	},
)

var RegisterSuccessTotal = promauto.NewCounter(
	prometheus.CounterOpts{
		Namespace: namespace,
		Subsystem: subsystem,
		Name:      "register_success_total",
		Help:      "Total number of successfull register requests",
	},
)

var RegisterErrorsWithReason = promauto.NewCounterVec(
	prometheus.CounterOpts{
		Namespace: namespace,
		Subsystem: subsystem,
		Name:      "register_errors_total",
		Help:      "Register error with reason",
	},
	[]string{"reason"},
)

var RegisterPasswordHashDurationSeconds = promauto.NewHistogram(
	prometheus.HistogramOpts{
		Namespace: namespace,
		Subsystem: subsystem,
		Name:      "register_password_hash_duration_seconds",
		Help:      "Password hash latency",
		Buckets:   []float64{0.01, 0.05, 0.1, 0.25, 0.5, 1, 2},
	},
)

var RegisterDBDurationSeconds = promauto.NewHistogram(
	prometheus.HistogramOpts{
		Namespace: namespace,
		Subsystem: subsystem,
		Name:      "register_db_duration_seconds",
		Help:      "Database connection latency",
		Buckets:   []float64{0.001, 0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1},
	},
)
