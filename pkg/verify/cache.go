/*
 * Copyright (c) 2019 vChain, Inc. All Rights Reserved.
 * This software is released under GPL3.
 * The full license information can be found under:
 * https://www.gnu.org/licenses/gpl-3.0.en.html
 *
 */
package verify

import (
	"sync"
)

type digestCache struct {
	resMap map[string]string
	mu     *sync.RWMutex
}

// Get a cached digest by key
func (d digestCache) Get(key string) string {
	d.mu.RLock()
	defer d.mu.RUnlock()
	return d.resMap[key]
}

// Set a cached digest by key
func (d digestCache) Set(key string, digest string) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.resMap[key] = digest
}

var dCache *digestCache

func init() {
	dCache = &digestCache{
		resMap: make(map[string]string),
		mu:     &sync.RWMutex{},
	}
}
