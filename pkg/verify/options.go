/*
 * Copyright (c) 2019 vChain, Inc. All Rights Reserved.
 * This software is released under GPL3.
 * The full license information can be found under:
 * https://www.gnu.org/licenses/gpl-3.0.en.html
 *
 */

package verify

import (
	"github.com/google/go-containerregistry/pkg/authn"
)

// Option is a functional option for verification operations
type Option func(*options) error

type options struct {
	keys     []string
	org      string
	keychain authn.Keychain
}

func makeOptions(opts ...Option) (*options, error) {
	o := &options{}

	for _, option := range opts {
		if err := option(o); err != nil {
			return nil, err
		}
	}
	return o, nil
}

func WithAuthKeychain(keychain authn.Keychain) Option {
	return func(o *options) error {
		o.keychain = keychain
		return nil
	}
}

func WithSignerKeys(keys ...string) Option {
	return func(o *options) error {
		o.keys = keys
		return nil
	}
}

func WithSignerOrg(org string) Option {
	return func(o *options) error {
		o.org = org
		return nil
	}
}
