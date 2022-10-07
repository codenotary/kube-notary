/*
 * Copyright (c) 2019 vChain, Inc. All Rights Reserved.
 * This software is released under GPL3.
 * The full license information can be found under:
 * https://www.gnu.org/licenses/gpl-3.0.en.html
 *
 */

package status

import (
	"embed"
	"net/http"
)

//go:embed status
var static embed.FS

// Handler returns an http.Handler to expose the status page.
func Handler() http.Handler {
	return http.StripPrefix("/status/", http.FileServer(http.FS(static)))
}
