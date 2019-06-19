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
	"strings"

	"github.com/google/go-containerregistry/pkg/authn"
	"github.com/google/go-containerregistry/pkg/name"
	containerregistry "github.com/google/go-containerregistry/pkg/v1"
	"github.com/google/go-containerregistry/pkg/v1/remote"
)

const (
	// DockerImageIDPrefix is the prefix of image id in container status.
	DockerImageIDPrefix = "docker://"
	// DockerPullableImageIDPrefix is the prefix of pullable image id in container status.
	DockerPullableImageIDPrefix = "docker-pullable://"
	// SHA256DigestPrefix is the prefix of supported image digest.
	SHA256DigestPrefix = "sha256:"
)

func image(imageRef string, keychain authn.Keychain) (containerregistry.Image, error) {
	ref, err := name.ParseReference(imageRef)
	if err != nil {
		return nil, fmt.Errorf("parsing reference %q: %v", imageRef, err)
	}

	return remote.Image(ref, remote.WithAuthFromKeychain(keychain))
}

func configDigest(imageRef string, keychain authn.Keychain) (string, error) {
	img, err := image(imageRef, keychain)
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
//
// See https://github.com/google/go-containerregistry/blob/master/images/ociimage.jpeg
// for an overiew about image ids and references.
//
// Supported formats:
//  - Local OCI image id (pre-pulled, `sha256:` prefix)
//  - Local docker image id (pre-pulled, `docker://sha256:` prefix)
//  - OCI pullable image reference
//  - Docker pullable image reference (prefix `docker-pullable://`)
//
// Resolution strategy:
//  - Pre-pulled images are resolved directly (digest is the actual image id already)
//  - For pullable image references, the manifest is fetched from the respective registry
//    using the auth keychain if needed. If manifest was a list, the image matching the
//    current platform will be chosen. Finally, the config digest is returned.
//
// Note:
//  - only sha256 digests are supported
//  - Docker Manifest v1 is not yet supported, see:
//    https://github.com/google/go-containerregistry/blob/master/pkg/v1/remote/descriptor.go#L111
//    https://github.com/google/go-containerregistry/issues/377
func Resolve(imageID string, keychain authn.Keychain) (string, error) {

	// Docker pre-pulled image
	if strings.HasPrefix(imageID, DockerImageIDPrefix+SHA256DigestPrefix) {
		return strings.TrimPrefix(imageID, DockerImageIDPrefix), nil
	}
	// OCI pre-pulled image
	if strings.HasPrefix(imageID, SHA256DigestPrefix) {
		return imageID, nil
	}

	// Pullable image
	if strings.Contains(imageID, "@"+SHA256DigestPrefix) {
		ref := imageID
		if strings.HasPrefix(imageID, DockerPullableImageIDPrefix) {
			ref = strings.TrimPrefix(imageID, DockerPullableImageIDPrefix)
		}
		return configDigest(ref, keychain)
	}

	return "", fmt.Errorf("unsupported image format: %s", imageID)
}
