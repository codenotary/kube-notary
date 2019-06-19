/*
 * Copyright (c) 2019 vChain, Inc. All Rights Reserved.
 * This software is released under GPL3.
 * The full license information can be found under:
 * https://www.gnu.org/licenses/gpl-3.0.en.html
 *
 */
package image

import (
	"fmt"
	"testing"
)

func TestResolve(t *testing.T) {
	testCases := map[string]string{
		"sha256:53f3fd8007f76bd23bf663ad5f5009c8941f63828ae458cef584b5f85dc0a7bf":                                           "sha256:53f3fd8007f76bd23bf663ad5f5009c8941f63828ae458cef584b5f85dc0a7bf",
		"docker://sha256:53f3fd8007f76bd23bf663ad5f5009c8941f63828ae458cef584b5f85dc0a7bf":                                  "sha256:53f3fd8007f76bd23bf663ad5f5009c8941f63828ae458cef584b5f85dc0a7bf",
		"docker.io/library/nginx@sha256:23b4dcdf0d34d4a129755fc6f52e1c6e23bb34ea011b315d87e193033bcd1b68":                   "sha256:53f3fd8007f76bd23bf663ad5f5009c8941f63828ae458cef584b5f85dc0a7bf",
		"docker-pullable://docker.io/library/nginx@sha256:23b4dcdf0d34d4a129755fc6f52e1c6e23bb34ea011b315d87e193033bcd1b68": "sha256:53f3fd8007f76bd23bf663ad5f5009c8941f63828ae458cef584b5f85dc0a7bf",
	}

	for in, expected := range testCases {
		out, err := Resolve(in, nil)
		if err != nil {
			t.Error(err)
		}
		if out != expected {
			t.Error(
				fmt.Sprintf(`Expected "%s", got "%s`, expected, out),
			)
		}
	}

	out, err := Resolve("invalid reference", nil)
	if err == nil {
		t.Error("Error expected")
	}
	if out != "" {
		t.Error("Empty string expected")
	}
}
