/*
 * Copyright (c) 2019 vChain, Inc. All Rights Reserved.
 * This software is released under GPL3.
 * The full license information can be found under:
 * https://www.gnu.org/licenses/gpl-3.0.en.html
 *
 */
package metrics

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	labelNames = []string{
		"namespace",
		"pod",
		"container",
		// "container_id", // not actually needed
		"image",
		"image_id",
		"hash",
	}

	verificationStatus = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "vcn_verification_status",
			Help: "Current verification status of images.",
		},
		labelNames,
	)

	verificationLevel = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "vcn_verification_level",
			Help: "Current verification level of images.",
		},
		labelNames,
	)
)

func init() {
	// Metrics have to be registered to be exposed:
	prometheus.MustRegister(verificationStatus)
	prometheus.MustRegister(verificationLevel)
}

func Handler() http.Handler {
	return promhttp.Handler()
}

func (p prometheusRecorder) Record(m Metric) {
	labels := prometheus.Labels{
		"namespace": m.Pod.Namespace,
		"pod":       m.Pod.Name,
		"container": m.ContainerStatus.Name,
		// "container_id": m.ContainerStatus.ContainerID,
		"image":    m.ContainerStatus.Image,
		"image_id": m.ContainerStatus.ImageID,
		"hash":     m.Hash,
	}

	verificationStatus.With(labels).Add(float64(m.Verification.Status))
	verificationLevel.With(labels).Add(float64(m.Verification.Level))
}

func (p prometheusRecorder) Reset() {
	verificationStatus.Reset()
	verificationLevel.Reset()
}

type prometheusRecorder struct{}

func NewRecorder() Recorder {
	return &prometheusRecorder{}
}
