package source

import (
	"github.com/jimeh/build-emacs-for-macos/pkg/commit"
	"github.com/jimeh/build-emacs-for-macos/pkg/repository"
)

type Source struct {
	Ref        string                 `yaml:"ref,omitempty" json:"ref,omitempty"`
	Repository *repository.Repository `yaml:"repository,omitempty" json:"repository,omitempty"`
	Commit     *commit.Commit         `yaml:"commit,omitempty" json:"commit,omitempty"`
	Tarball    *Tarball               `yaml:"tarball,omitempty" json:"tarball,omitempty"`
}

type Tarball struct {
	URL string `yaml:"url,omitempty" json:"url,omitempty"`
}
