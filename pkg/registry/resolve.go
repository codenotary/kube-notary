/*
 * Copyright (c) 2019 vChain, Inc. All Rights Reserved.
 * This software is released under GPL3.
 * The full license information can be found under:
 * https://www.gnu.org/licenses/gpl-3.0.en.html
 *
 */
package registry

import (
	"fmt"
	"strings"

	"github.com/google/go-containerregistry/pkg/authn"
	"github.com/google/go-containerregistry/pkg/name"
	containerregistry "github.com/google/go-containerregistry/pkg/v1"
	"github.com/google/go-containerregistry/pkg/v1/remote"
)

// Image provides access to an image reference
func Image(imageRef string, keychain authn.Keychain) (containerregistry.Image, error) {
	ref, err := name.ParseReference(imageRef)
	if err != nil {
		return nil, fmt.Errorf("parsing reference %q: %v", imageRef, err)
	}

	return remote.Image(ref, remote.WithAuthFromKeychain(keychain))
}

func imageConfigDigest(imageRef string, keychain authn.Keychain) (string, error) {
	img, err := Image(imageRef, keychain)
	if err != nil {
		return "", fmt.Errorf("reading image %s: %v", imageRef, err)
	}

	configName, err := img.ConfigName()
	if err != nil {
		return "", err
	}

	return configName.String(), nil
}

// Resolve returns the actual image id (ie. the digest of the image's configuration)
// from a given ImageID of the container's image as per Kubernetes specs.
// See https://github.com/google/go-containerregistry/blob/master/images/ociimage.jpeg
func Resolve(imageID string, keychain authn.Keychain) (string, error) {

	if strings.HasPrefix(imageID, "sha256:") { // not-pullable image
		return imageID, nil
	}

	if strings.Contains(imageID, "@sha256:") { // pullable image
		var ref string
		switch true {
		// see https://github.com/kubernetes/kubernetes/blob/master/pkg/kubelet/dockershim/convert.go#L68
		case strings.HasPrefix(imageID, "docker-pullable://"):
			ref = strings.TrimPrefix(imageID, "docker-pullable://")
		// containerd
		case !strings.Contains(imageID, "://"):
			ref = imageID
		default:
			return "", fmt.Errorf("unsupported image format: %s", imageID)
		}
		return imageConfigDigest(ref, keychain)
	}

	return "", fmt.Errorf("unsupported image format: %s", imageID)
}
