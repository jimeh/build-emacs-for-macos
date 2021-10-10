package osinfo

import (
	"os/exec"
	"strconv"
	"strings"
)

type OSInfo struct {
	Name    string `yaml:"name" json:"name"`
	Version string `yaml:"version" json:"version"`
	Arch    string `yaml:"arch" json:"arch"`
}

func New() (*OSInfo, error) {
	version, err := exec.Command("sw_vers", "-productVersion").CombinedOutput()
	if err != nil {
		return nil, err
	}

	arch, err := exec.Command("uname", "-m").CombinedOutput()
	if err != nil {
		return nil, err
	}

	return &OSInfo{
		Name:    "macOS",
		Version: strings.TrimSpace(string(version)),
		Arch:    strings.TrimSpace(string(arch)),
	}, nil
}

// DistinctVersion returns macOS version down to a distinct "major"
// version. For macOS 10.x, this will include the first two numeric parts of the
// version (10.15), while for 11.x and later, the first numeric part is enough
// (11).
func (s *OSInfo) DistinctVersion() string {
	parts := strings.Split(s.Version, ".")

	if n, _ := strconv.Atoi(parts[0]); n >= 11 {
		return parts[0]
	}

	max := len(parts)
	if max > 2 {
		max = 2
	}

	return strings.Join(parts[0:max], ".")
}
