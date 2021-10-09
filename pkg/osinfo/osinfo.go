package osinfo

import (
	"os/exec"
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

func (s *OSInfo) MajorMinor() string {
	parts := strings.Split(s.Version, ".")
	max := len(parts)
	if max > 2 {
		max = 2
	}

	return strings.Join(parts[0:max], ".")
}
