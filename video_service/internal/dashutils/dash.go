// This package provides utilities for transcoding/chunking in compliance with DASH's requirements
package dashutils

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"os/exec"
	"path/filepath"
)

type DASHVideo struct {
	ManifestPath string
	QualityMap   []string
}

func TranscodeAndGenerateManifest(path string, local bool) (*DASHVideo, error) {
	var encodeArgs []string
	switch local {
	case true:
		// make the encoding really fast so we don't have to wait 90 minutes for integration tests
		// I don't really understand the difference between the  speed and deadline args, but the documentation implies
		// they're separate
		encodeArgs = []string{path, "-speed 16 -deadline realtime -r 1 -crf 63 -t 10"}
	case false:
		encodeArgs = []string{path, "-speed 3 -deadline good -crf 30"}
	}

	cmd := exec.Command("./scripts/transcode.sh", encodeArgs...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Errorf("%s", out)
		return nil, err
	}

	// At this point it's been transcoded, so generate the DASH manifest
	cmd = exec.Command("./scripts/manifest.sh", path)
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
