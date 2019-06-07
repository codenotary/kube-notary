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
	"github.com/google/go-containerregistry/pkg/v1/remote"
)

func resolveManifest(image string) (string, error) {
	ref, err := name.ParseReference(image)
	if err != nil {
		return "", fmt.Errorf("parsing reference %q: %v", image, err)
	}

	img, err := remote.Image(ref, remote.WithAuthFromKeychain(authn.DefaultKeychain))
	if err != nil {
		return "", fmt.Errorf("reading image %q: %v", ref, err)
	}

	configName, err := img.ConfigName()
	if err != nil {
		return "", err
	}

	return configName.String(), nil
}

// Digest returns the actual image's digest as string
func Digest(imageID string) (string, error) {

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
		return resolveManifest(ref)
	}

	return "", fmt.Errorf("unsupported image format: %s", imageID)
}
