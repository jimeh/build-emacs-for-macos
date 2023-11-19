package cask

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAsset(t *testing.T) {
	// Define test cases
	tests := []struct {
		name    string
		release ReleaseInfo
		needles []string
		want    *ReleaseAsset
	}{
		{
			name: "single needle, exact match",
			release: ReleaseInfo{
				Assets: map[string]*ReleaseAsset{
					"asset1": {Filename: "asset1.zip"},
					"asset2": {Filename: "asset2.zip"},
				},
			},
			needles: []string{"asset1"},
			want:    &ReleaseAsset{Filename: "asset1.zip"},
		},
		{
			name: "multiple needles, all",
			release: ReleaseInfo{
				Assets: map[string]*ReleaseAsset{
					"asset1": {Filename: "asset1.zip"},
					"asset2": {Filename: "asset2.zip"},
				},
			},
			needles: []string{"zip", "asset1"},
			want:    &ReleaseAsset{Filename: "asset1.zip"},
		},
		{
			name: "multiple needles, one match",
			release: ReleaseInfo{
				Assets: map[string]*ReleaseAsset{
					"asset1": {Filename: "asset1.zip"},
					"asset2": {Filename: "asset2.zip"},
				},
			},
			needles: []string{"rar", "asset2"},
			want:    nil,
		},
		{
			name: "multiple needles, no match",
			release: ReleaseInfo{
				Assets: map[string]*ReleaseAsset{
					"asset1": {Filename: "asset1.zip"},
					"asset2": {Filename: "asset2.zip"},
				},
			},
			needles: []string{"rar", "asset3"},
			want:    nil,
		},
		{
			name: "no needles",
			release: ReleaseInfo{
				Assets: map[string]*ReleaseAsset{
					"asset1": {Filename: "asset1.zip"},
					"asset2": {Filename: "asset2.zip"},
				},
			},
			needles: nil,
			want:    &ReleaseAsset{Filename: "asset1.zip"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.release.Asset(tt.needles...)

			assert.Equal(t, tt.want, got)
		})
	}
}
