package plan

import (
	"bytes"
	"io"
	"os"

	"github.com/jimeh/build-emacs-for-macos/pkg/osinfo"
	"github.com/jimeh/build-emacs-for-macos/pkg/source"
	"gopkg.in/yaml.v3"
)

type Plan struct {
	Build   *Build         `yaml:"build,omitempty"`
	Source  *source.Source `yaml:"source,omitempty"`
	OS      *osinfo.OSInfo `yaml:"os,omitempty"`
	Release *Release       `yaml:"release,omitempty"`
	Output  *Output        `yaml:"output,omitempty"`
}

// Load attempts to loads a plan YAML from given filename.
func Load(filename string) (*Plan, error) {
	b, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	p := &Plan{}
	err = yaml.Unmarshal(b, p)
	if err != nil {
		return nil, err
	}

	return p, nil
}

// WriteYAML writes plan in YAML format to given io.Writer.
func (s *Plan) WriteYAML(w io.Writer) error {
	enc := yaml.NewEncoder(w)
	enc.SetIndent(2)

	return enc.Encode(s)
}

// YAML returns plan in YAML format.
func (s *Plan) YAML() (string, error) {
	var buf bytes.Buffer
	err := s.WriteYAML(&buf)
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}

type Build struct {
	Name string `yaml:"name,omitempty"`
}

type Release struct {
	Name       string `yaml:"name"`
	Title      string `yaml:"title,omitempty"`
	Draft      bool   `yaml:"draft,omitempty"`
	Prerelease bool   `yaml:"prerelease,omitempty"`
}

type Output struct {
	Directory string `yaml:"directory,omitempty"`
	DiskImage string `yaml:"disk_image,omitempty"`
}
