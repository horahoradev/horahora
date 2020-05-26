// This package provides utilities for transcoding/chunking in compliance with DASH's requirements
package dashutils

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"

	log "github.com/sirupsen/logrus"
)

type DASHVideo struct {
	ManifestPath string
	QualityMap   []string
}

func TranscodeAndChunk(path string) (*DASHVideo, error) {
	// FIXME: remove extension from path to maintain compatibility with bad scripts
	// I should fix this!
	path = strings.Replace(path, ".mp4", "", 1)

	cmd := exec.Command("./scripts/transcode.sh", path)
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Errorf("%s", out)
		return nil, err
	}

	// At this point it's been transcoded, so chunk/generate the DASH manifest
	cmd = exec.Command("./scripts/chunk.sh", path)
	_, err = cmd.CombinedOutput()
	if err != nil {
		return nil, err
	}

	var fileList []string

	generatedFiles, err := filepath.Glob(fmt.Sprintf("%s_*", path))
	if err != nil {
		return nil, err
	}

	for _, fileName := range generatedFiles {
		fileList = append(fileList, fileName)
	}

	return &DASHVideo{
		ManifestPath: fmt.Sprintf("%s.mpd", path),
		QualityMap:   fileList,
	}, nil
}
