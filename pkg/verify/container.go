/*
 * Copyright (c) 2019 vChain, Inc. All Rights Reserved.
 * This software is released under GPL3.
 * The full license information can be found under:
 * https://www.gnu.org/licenses/gpl-3.0.en.html
 *
 */

package verify

import (
	"github.com/vchain-us/kube-notary/pkg/image"
	"github.com/vchain-us/vcn/pkg/api"
	corev1 "k8s.io/api/core/v1"
)

// ContainerStatus returns the hast string and the BlockchainVerification for the given status
func ContainerStatus(status corev1.ContainerStatus, options ...Option) (hash string, verification *api.BlockchainVerification, err error) {
	o, err := makeOptions(options...)
	if err != nil {
		return
	}

	// Cache image to digest on container basis
	cKey := status.ContainerID + "|" + status.ImageID
	digest := dCache.Get(cKey)
	if digest == "" {
		digest, err = image.Resolve(status.ImageID, o.keychain)
		if err != nil {
			return
		}
		dCache.Set(cKey, digest)
	}

	return getVerification(digest, o)
}
