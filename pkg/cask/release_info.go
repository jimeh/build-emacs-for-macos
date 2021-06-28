package cask

import (
	"sort"
	"strings"
)

type ReleaseInfo struct {
	Name    string
	Version string
	Assets  map[string]*ReleaseAsset
}

func (s *ReleaseInfo) Asset(nameMatch string) *ReleaseAsset {
	if a, ok := s.Assets[nameMatch]; ok {
		return a
	}

	// Dirty and inefficient way to ensure assets are searched in a predictable
	// order.
	var assets []*ReleaseAsset
	for _, a := range s.Assets {
		assets = append(assets, a)
	}
	sort.SliceStable(assets, func(i, j int) bool {
		return assets[i].Filename < assets[j].Filename
	})

	for _, a := range assets {
		if strings.Contains(a.Filename, nameMatch) {
			return a
		}
	}

	return nil
}

func (s *ReleaseInfo) DownloadURL(nameMatch string) string {
	a := s.Asset(nameMatch)
	if a == nil {
		return ""
	}

	return a.DownloadURL
}

func (s *ReleaseInfo) SHA256(nameMatch string) string {
	a := s.Asset(nameMatch)
	if a == nil {
		return ""
	}

	return a.SHA256
}

type ReleaseAsset struct {
	Filename    string
	DownloadURL string
	SHA256      string
}
