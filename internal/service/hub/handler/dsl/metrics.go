package dsl

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	requestCounter = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "dsl_get_activity_requests_total",
			Help: "Total number of GetActivity requests",
		},
	)
	requestCounterByDirection = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "dsl_get_activity_requests_by_direction_total",
			Help: "Total number of GetActivity requests by direction",
		},
		[]string{"direction"},
	)
	requestCounterByNetwork = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "dsl_get_activity_requests_by_network_total",
			Help: "Total number of GetActivity requests by network",
		},
		[]string{"network"},
	)
	requestCounterByTag = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "dsl_get_activity_requests_by_tag_total",
			Help: "Total number of GetActivity requests by tag",
		},
		[]string{"tag"},
	)
	requestCounterByType = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "dsl_get_activity_requests_by_type_total",
			Help: "Total number of GetActivity requests by type",
		},
		[]string{"type"},
	)
	requestCounterByPlatform = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "dsl_get_activity_requests_by_platform_total",
			Help: "Total number of GetActivity requests by platform",
		},
		[]string{"platform"},
	)
)

func incrementRequestCounter(direction *string, network []string, tag []string, platform []string, theType []string) {
	requestCounter.Inc()

	if direction != nil {
		requestCounterByDirection.WithLabelValues(*direction).Inc()
	}
	for _, t := range network {
		requestCounterByNetwork.WithLabelValues(t).Inc()
	}
	for _, t := range tag {
		requestCounterByTag.WithLabelValues(t).Inc()
	}
	for _, t := range platform {
		requestCounterByPlatform.WithLabelValues(t).Inc()
	}
	for _, t := range theType {
		requestCounterByType.WithLabelValues(t).Inc()
	}
}
