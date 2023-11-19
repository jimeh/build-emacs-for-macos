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

func (s *ReleaseInfo) Asset(needles ...string) *ReleaseAsset {
	if len(needles) == 1 {
		if a, ok := s.Assets[needles[0]]; ok {
			return a
		}
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

assetsLoop:
	for _, a := range assets {
		for _, needle := range needles {
			if !strings.Contains(a.Filename, needle) {
				continue assetsLoop
			}
		}

		return a
	}

	return nil
}

func (s *ReleaseInfo) DownloadURL(needles ...string) string {
	a := s.Asset(needles...)
	if a == nil {
		return ""
	}

	return a.DownloadURL
}

func (s *ReleaseInfo) SHA256(needles ...string) string {
	a := s.Asset(needles...)
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
