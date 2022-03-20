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
)

// ImageHash returns the hash string and the verification status for the given imageID
func ImageHash(imageID string, options ...Option) (hash string, err error) {
	o, err := makeOptions(options...)
	if err != nil {
		return
	}
	return image.Resolve(imageID, o.keychain)
}
