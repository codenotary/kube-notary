/*
 * Copyright (c) 2019 vChain, Inc. All Rights Reserved.
 * This software is released under GPL3.
 * The full license information can be found under:
 * https://www.gnu.org/licenses/gpl-3.0.en.html
 *
 */
package metrics

import (
	"github.com/vchain-us/vcn/pkg/api"

	log "github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
)

type Metric struct {
	Pod             *corev1.Pod
	ContainerStatus *corev1.ContainerStatus
	Verification    *api.BlockchainVerification
	Hash            string
}

func (m Metric) LogFields() *log.Fields {
	return &log.Fields{
		"namespace":           m.Pod.Namespace,
		"pod":                 m.Pod.Name,
		"container":           m.ContainerStatus.Name,
		"container_id":        m.ContainerStatus.ContainerID,
		"image":               m.ContainerStatus.Image,
		"image_id":            m.ContainerStatus.ImageID,
		"hash":                m.Hash,
		"verification_date":   m.Verification.Date(),
		"verification_level":  m.Verification.Level,
		"verification_status": m.Verification.Status,
		"trusted":             m.Verification.Trusted(),
	}
}
