/*
 * Copyright (c) 2019 vChain, Inc. All Rights Reserved.
 * This software is released under GPL3.
 * The full license information can be found under:
 * https://www.gnu.org/licenses/gpl-3.0.en.html
 *
 */
package registry

import (
	"github.com/google/go-containerregistry/pkg/authn"
	"github.com/google/go-containerregistry/pkg/authn/k8schain"
	"k8s.io/client-go/kubernetes"
)

// NewKeychain returns a new authn.Keychain suitable for resolving image references as
// scoped by the provided namespace, serviceAccountName, and imagePullSecretes.
// It speaks to Kubernetes through the provided client interface.
func NewKeychain(client kubernetes.Interface, namespace string, serviceAccountName string, imagePullSecrets []string) (authn.Keychain, error) {
	return k8schain.New(client, k8schain.Options{
		Namespace:          namespace,
		ServiceAccountName: serviceAccountName,
		ImagePullSecrets:   imagePullSecrets,
	})
}
