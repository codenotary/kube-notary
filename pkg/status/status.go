/*
 * Copyright (c) 2019 vChain, Inc. All Rights Reserved.
 * This software is released under GPL3.
 * The full license information can be found under:
 * https://www.gnu.org/licenses/gpl-3.0.en.html
 *
 */
//go:generate statik -f -src=./status
package status

import (
	"net/http"

	"github.com/rakyll/statik/fs"

	// embedded static files
	_ "github.com/vchain-us/kube-notary/pkg/status/statik"
)

// Handler returns an http.Handler to expose the status page.
func Handler() http.Handler {
	statikFS, err := fs.New()
	if err != nil {
		panic(err)
	}
	return http.StripPrefix("/status/", http.FileServer(statikFS))
}
