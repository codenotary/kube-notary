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

func getVerification(digest string, o *options) (hash string, verification *api.BlockchainVerification, err error) {
	hash = strings.TrimPrefix(digest, "sha256:")
	if o.org != "" {
		bo, err := api.BlockChainGetOrganisation(o.org)
		if err != nil {
			return hash, nil, err
		}
		verification, err = api.BlockChainVerifyMatchingPublicKeys(hash, bo.MembersKeys())
	} else if len(o.keys) > 0 {
		verification, err = api.BlockChainVerifyMatchingPublicKeys(hash, o.keys)
	} else {
		verification, err = api.BlockChainVerify(hash)
	}
	api.TrackVerify(nil, hash, digest)
	return
}
