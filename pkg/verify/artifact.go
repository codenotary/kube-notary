/*
 * Copyright (c) 2019 vChain, Inc. All Rights Reserved.
 * This software is released under GPL3.
 * The full license information can be found under:
 * https://www.gnu.org/licenses/gpl-3.0.en.html
 *
 */

package verify

import (
	"github.com/vchain-us/vcn/pkg/api"
)

// Artifact gets and returns the api.ArtifactResponse from platform API, if any
func Artifact(hash string, verification *api.BlockchainVerification) (*api.ArtifactResponse, error) {
	if verification == nil {
		return nil, nil
	}
	return api.LoadArtifactForHash(nil, hash, verification.MetaHash())
}
