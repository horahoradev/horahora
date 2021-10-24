package dashutils

// This implementation just uses the original encoding

type NullTrancoder struct {}

func (h NullTrancoder) TranscodeAndGenerateManifest(path string, local bool) (*DASHVideo, error) {
	return &DASHVideo{
		ManifestPath:     nil,
		QualityMap:       []string{},
		OriginalFilePath: path,
		ThumbnailPath:    path + ".jpg",
	}, nil
}