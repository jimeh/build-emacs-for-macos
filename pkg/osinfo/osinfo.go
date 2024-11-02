package osinfo

import (
	"os"
	"os/exec"
	"strconv"
	"strings"
)

type OSInfo struct {
	Name       string `yaml:"name" json:"name"`
	Version    string `yaml:"version" json:"version"`
	SDKVersion string `yaml:"sdk_version" json:"sdk_version"`
	Arch       string `yaml:"arch" json:"arch"`
}

func New() (*OSInfo, error) {
	version, err := exec.Command("sw_vers", "-productVersion").Output()
	if err != nil {
		return nil, err
	}

	sdkVersion := os.Getenv("MACOSX_DEPLOYMENT_TARGET")
	if sdkVersion == "" {
		var ver []byte
		ver, err = exec.Command("xcrun", "--show-sdk-version").Output()
		if err != nil {
			return nil, err
		}

		sdkVersion = string(ver)
	}

	arch, err := exec.Command("uname", "-m").CombinedOutput()
	if err != nil {
		return nil, err
	}

	return &OSInfo{
		Name:       "macOS",
		Version:    strings.TrimSpace(string(version)),
		SDKVersion: strings.TrimSpace(sdkVersion),
		Arch:       strings.TrimSpace(string(arch)),
	}, nil
}

// DistinctVersion returns macOS version down to a distinct "major" version. For
// macOS 10.x, this will include the first two numeric parts of the version
// (10.15), while for 11.x and later, the first numeric part is enough (11).
func (s *OSInfo) DistinctVersion() string {
	return s.distinctVersion(s.Version)
}

// DistinctSDKVersion returns macOS version down to a distinct "major" version.
// For macOS 10.x, this will include the first two numeric parts of the version
// (10.15), while for 11.x and later, the first numeric part is enough (11).
func (s *OSInfo) DistinctSDKVersion() string {
	return s.distinctVersion(s.SDKVersion)
}

func (s *OSInfo) distinctVersion(version string) string {
	parts := strings.Split(version, ".")

	if n, _ := strconv.Atoi(parts[0]); n >= 11 {
		return parts[0]
	}

	max := len(parts)
	if max > 2 {
		max = 2
	}

	return strings.Join(parts[0:max], ".")
}
