/*
 * Copyright (c) 2019 vChain, Inc. All Rights Reserved.
 * This software is released under GPL3.
 * The full license information can be found under:
 * https://www.gnu.org/licenses/gpl-3.0.en.html
 *
 */
package image

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
)

type cache struct {
	resMap map[string]string
	mu     *sync.RWMutex
}

//Get a cached hash by imageID
func (c cache) Get(imageID string) string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.resMap[imageID]
}

// Set a cached hash by imageID
func (c cache) Set(imageID string, hash string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.resMap[imageID] = hash
}

// CacheHandler returns an http.Handler to expose the internal image resolution cache.
func CacheHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c := hashCache
		c.mu.RLock()
		defer c.mu.RUnlock()

		b, err := json.Marshal(c.resMap)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintln(w, err.Error())
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(b)
	})
}

var hashCache *cache

func init() {
	hashCache = &cache{
		resMap: make(map[string]string),
		mu:     &sync.RWMutex{},
	}
}
