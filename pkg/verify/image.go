/*
 * Copyright (c) 2019 vChain, Inc. All Rights Reserved.
 * This software is released under GPL3.
 * The full license information can be found under:
 * https://www.gnu.org/licenses/gpl-3.0.en.html
 *
 */

package verify

import (
	"strings"

	"github.com/vchain-us/vcn/pkg/api"
)

// ImageID returns the BlockchainVerification for imageId
func ImageID(imageID string, signerKeys ...string) (verification *api.BlockchainVerification, err error) {
	hash := strings.TrimPrefix(imageID, "sha256:")
	if len(signerKeys) > 0 {
		verification, err = api.BlockChainVerifyMatchingPublicKeys(hash, signerKeys)
	} else {
		verification, err = api.BlockChainVerify(hash)
	}
	return
}
