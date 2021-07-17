// This package provides utilities for transcoding/chunking in compliance with DASH's requirements
package dashutils

import (
	"fmt"
	"os/exec"
	"path/filepath"

	log "github.com/sirupsen/logrus"
)

type DASHVideo struct {
	ManifestPath     string
	ThumbnailPath    string
	QualityMap       []string
	OriginalFilePath string
}

func TranscodeAndGenerateManifest(path string, local bool) (*DASHVideo, error) {
	// var encodeArgs []string
	// switch local {
	// case true:
	// 	// make the encoding really fast so we don't have to wait 90 minutes for integration tests
	// 	// I don't really understand the difference between the speed and deadline args, but the documentation implies
	// 	// they're separate
	// 	encodeArgs = []string{path, "-speed 16 -deadline realtime -r 1 -crf 63 -t 10"}
	// case false:
	// 	// -r 24 -deadline realtime -cpu-used 1
	// 	encodeArgs = []string{path, "-r 24 -deadline good -cpu-used 2"}
	// }

	cmd := exec.Command("/videoservice/scripts/transcode.sh", []string{path}...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Errorf("%s", out)
		return nil, err
	}

	// At this point it's been transcoded, so generate the DASH manifest
	cmd = exec.Command("/videoservice/scripts/manifest.sh", []string{filepath.Base(path)}...)
	_, err = cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("failed to generate dash manifest. Err: %s", err)
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
		ManifestPath:     fmt.Sprintf("%s.mpd", path),
		QualityMap:       fileList,
		OriginalFilePath: path,
		ThumbnailPath:    path + ".jpg",
	}, nil
}
