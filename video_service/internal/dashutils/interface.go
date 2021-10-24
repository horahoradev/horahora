// This package provides utilities for transcoding/chunking in compliance with DASH's requirements
package dashutils

import (
)

type DASHVideo struct {
	ManifestPath     *string
	ThumbnailPath    string
	QualityMap       []string
	OriginalFilePath string
}

type Transcoder interface {
	TranscodeAndGenerateManifest(path string, local bool) (*DASHVideo, error)
}
