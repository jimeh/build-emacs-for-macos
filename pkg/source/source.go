package source

import (
	"github.com/jimeh/build-emacs-for-macos/pkg/commit"
	"github.com/jimeh/build-emacs-for-macos/pkg/repository"
)

type Source struct {
	Ref        string                 `yaml:"ref,omitempty"`
	Repository *repository.Repository `yaml:"repository,omitempty"`
	Commit     *commit.Commit         `yaml:"commit,omitempty"`
	Tarball    *Tarball               `yaml:"tarball,omitempty"`
}

type Tarball struct {
	URL string `yaml:"url,omitempty"`
}
