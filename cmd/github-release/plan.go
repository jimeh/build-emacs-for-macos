package main

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"gopkg.in/yaml.v3"
)

type Plan struct {
	Ref     string `yaml:"ref"`
	SHA     string `yaml:"sha"`
	Date    string `yaml:"date"`
	Archive string `yaml:"archive"`
}

func LoadPlan(filename string) (*Plan, error) {
	b, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	plan := &Plan{}
	err = yaml.Unmarshal(b, plan)

	return plan, err
}

func (s *Plan) ReleaseName() string {
	ref := nonAlphaNum.ReplaceAllString(s.Ref, "-")
	ref = regexp.MustCompile(`\.`).ReplaceAllString(ref, "-")
	if ref == "master" {
		ref = "nightly"
	}

	return fmt.Sprintf("Emacs.%s.%s.%s", s.Date, s.SHA[0:6], ref)
}

func (s *Plan) ReleaseAsset() string {
	name := filepath.Base(s.Archive)
	ext := filepath.Ext(s.Archive)

	name = name[0 : len(name)-len(ext)]
	name = regexp.MustCompile(`^Emacs\.app-`).ReplaceAllString(name, "Emacs")
	name = regexp.MustCompile(`\.`).ReplaceAllString(name, "-")
	name = nonAlphaNum.ReplaceAllString(name, ".")
	name = strings.TrimRight(name, ".")

	return name + ext
}
