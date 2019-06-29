/*
 * Copyright (c) 2019 vChain, Inc. All Rights Reserved.
 * This software is released under GPL3.
 * The full license information can be found under:
 * https://www.gnu.org/licenses/gpl-3.0.en.html
 *
 */

package watcher

import (
	"io"

	"github.com/vchain-us/kube-notary/internal/util"
)

func bulkSigningScript(writer io.Writer, results []Result) error {
	return util.BulkSigningScriptTemplate.Execute(writer, struct {
		Results []Result
	}{Results: results})
}
