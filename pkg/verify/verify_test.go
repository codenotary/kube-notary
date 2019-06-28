/*
* Copyright (c) 2019 vChain, Inc. All Rights Reserved.
* This software is released under GPL3.
* The full license information can be found under:
* https://www.gnu.org/licenses/gpl-3.0.en.html
*
 */

package verify

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetVerification(t *testing.T) {
	// hash "test" was signed by 0x7f66cb537c27251d007bd3c8ec731690c744f5e4 onto blockchain for testing purpose

	h, v, e := getVerification("sha256:test", &options{})
	assert.Equal(t, "test", h)
	assert.NotNil(t, v)
	assert.NoError(t, e)

	h, v, e = getVerification("sha256:test", &options{keys: []string{"0x7f66cb537c27251d007bd3c8ec731690c744f5e4"}})
	assert.Equal(t, "test", h)
	assert.NotNil(t, v)
	assert.NoError(t, e)

}
