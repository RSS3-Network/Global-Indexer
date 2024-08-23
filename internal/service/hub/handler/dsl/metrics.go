package dsl

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/rss3-network/node/schema/worker/decentralized"
	"github.com/rss3-network/protocol-go/schema/network"
	"github.com/rss3-network/protocol-go/schema/tag"
)

var (
	requestCounter = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "dsl_get_activity_requests_total",
			Help: "Total number of GetActivity requests",
		},
		[]string{"endpoint"},
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
	requestCounterByPlatform = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "dsl_get_activity_requests_by_platform_total",
			Help: "Total number of GetActivity requests by platform",
		},
		[]string{"endpoint", "platform"},
	)
)

func incrementRequestCounter(endpoint string, networks []string, tags []string, platforms []string) {
	requestCounter.WithLabelValues(endpoint).Inc()

	if len(networks) > 0 {
		for _, t := range networks {
			requestCounterByNetwork.WithLabelValues(endpoint, t).Inc()
		}
	} else {
		for _, item := range network.NetworkStrings() {
			requestCounterByNetwork.WithLabelValues(endpoint, item).Inc()
		}
	}

	if len(tags) > 0 {
		for _, t := range tags {
			requestCounterByTag.WithLabelValues(endpoint, t).Inc()
		}
	} else {
		for _, item := range tag.TagStrings() {
			requestCounterByTag.WithLabelValues(endpoint, item).Inc()
		}
	}

	if len(platforms) > 0 {
		for _, t := range platforms {
			requestCounterByPlatform.WithLabelValues(endpoint, t).Inc()
		}
	} else {
		for _, item := range decentralized.PlatformStrings() {
			requestCounterByPlatform.WithLabelValues(endpoint, item).Inc()
		}
	}
}
