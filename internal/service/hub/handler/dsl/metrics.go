package dsl

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	requestCounter = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "dsl_get_activity_requests_total",
			Help: "Total number of GetActivity requests",
		},
		[]string{"endpoint"},
	)
	requestCounterByDirection = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "dsl_get_activity_requests_by_direction_total",
			Help: "Total number of GetActivity requests by direction",
		},
		[]string{"endpoint", "direction"},
	)
	requestCounterByNetwork = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "dsl_get_activity_requests_by_network_total",
			Help: "Total number of GetActivity requests by network",
		},
		[]string{"endpoint", "network"},
	)
	requestCounterByTag = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "dsl_get_activity_requests_by_tag_total",
			Help: "Total number of GetActivity requests by tag",
		},
		[]string{"endpoint", "tag"},
	)
	requestCounterByType = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "dsl_get_activity_requests_by_type_total",
			Help: "Total number of GetActivity requests by type",
		},
		[]string{"endpoint", "type"},
	)
	requestCounterByPlatform = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "dsl_get_activity_requests_by_platform_total",
			Help: "Total number of GetActivity requests by platform",
		},
		[]string{"endpoint", "platform"},
	)
)

func incrementRequestCounter(endpoint string, direction *string, network []string, tag []string, platform []string, theType []string) {
	requestCounter.WithLabelValues(endpoint).Inc()

	if direction != nil {
		requestCounterByDirection.WithLabelValues(endpoint, *direction).Inc()
	}

	for _, t := range network {
		requestCounterByNetwork.WithLabelValues(endpoint, t).Inc()
	}

	for _, t := range tag {
		requestCounterByTag.WithLabelValues(endpoint, t).Inc()
	}

	for _, t := range platform {
		requestCounterByPlatform.WithLabelValues(endpoint, t).Inc()
	}

	for _, t := range theType {
		requestCounterByType.WithLabelValues(endpoint, t).Inc()
	}
}
