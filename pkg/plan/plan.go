package plan

import (
	"bytes"
	"encoding/json"
	"io"
	"os"

	"github.com/jimeh/build-emacs-for-macos/pkg/osinfo"
	"github.com/jimeh/build-emacs-for-macos/pkg/source"
	"gopkg.in/yaml.v3"
)

type Plan struct {
	Build   *Build         `yaml:"build,omitempty" json:"build,omitempty"`
	Source  *source.Source `yaml:"source,omitempty" json:"source,omitempty"`
	OS      *osinfo.OSInfo `yaml:"os,omitempty" json:"os,omitempty"`
	Release *Release       `yaml:"release,omitempty" json:"release,omitempty"`
	Output  *Output        `yaml:"output,omitempty" json:"output,omitempty"`
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

// WriteJSON writes plan in JSON format to given io.Writer.
func (s *Plan) WriteJSON(w io.Writer) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")

	return enc.Encode(s)
}

// JSON returns plan in JSON format.
func (s *Plan) JSON() (string, error) {
	var buf bytes.Buffer
	err := s.WriteJSON(&buf)
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}

type Build struct {
	Name string `yaml:"name,omitempty" json:"name,omitempty"`
}

type Release struct {
	Name       string `yaml:"name" json:"name"`
	Title      string `yaml:"title,omitempty" json:"title,omitempty"`
	Draft      bool   `yaml:"draft,omitempty" json:"draft,omitempty"`
	Prerelease bool   `yaml:"prerelease,omitempty" json:"prerelease,omitempty"`
}

type Output struct {
	Directory string `yaml:"directory,omitempty" json:"directory,omitempty"`
	DiskImage string `yaml:"disk_image,omitempty" json:"disk_image,omitempty"`
}
