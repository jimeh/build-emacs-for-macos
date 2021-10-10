package repository

import (
	"errors"
	"fmt"
	"strings"
)

//nolint:golint
var (
	Err       = errors.New("repository")
	ErrGitHub = fmt.Errorf("%w: github", Err)
)

const GitHubBaseURL = "https://github.com/"

// Type is a repository type
type Type string

const GitHub Type = "github"

// Repository represents basic information about a repository with helper
// methods to get various pieces of information from it.
type Repository struct {
	Type   Type   `yaml:"type,omitempty" json:"type,omitempty"`
	Source string `yaml:"source,omitempty" json:"source,omitempty"`
}

func NewGitHub(ownerAndName string) (*Repository, error) {
	parts := strings.Split(ownerAndName, "/")
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return nil, fmt.Errorf(
			"%w: repository must be give in \"owner/name\" format",
			ErrGitHub,
		)
	}

	return &Repository{
		Type:   GitHub,
		Source: ownerAndName,
	}, nil
}

func (s *Repository) Owner() string {
	switch s.Type {
	case GitHub:
		return strings.SplitN(s.Source, "/", 2)[0]
	default:
		return ""
	}
}

func (s *Repository) Name() string {
	switch s.Type {
	case GitHub:
		return strings.SplitN(s.Source, "/", 2)[1]
	default:
		return ""
	}
}

func (s *Repository) URL() string {
	switch s.Type {
	case GitHub:
		return GitHubBaseURL + s.Source
	default:
		return ""
	}
}

func (s *Repository) CloneURL() string {
	switch s.Type {
	case GitHub:
		return GitHubBaseURL + s.Source + ".git"
	default:
		return ""
	}
}

func (s *Repository) TarballURL(ref string) string {
	if ref == "" {
		return ""
	}

	switch s.Type {
	case GitHub:
		return GitHubBaseURL + s.Source + "/tarball/" + ref
	default:
		return ""
	}
}

func (s *Repository) CommitURL(ref string) string {
	if ref == "" {
		return ""
	}

	switch s.Type {
	case GitHub:
		return GitHubBaseURL + s.Source + "/commit/" + ref
	default:
		return ""
	}
}

func (s *Repository) TreeURL(ref string) string {
	if ref == "" {
		return ""
	}

	switch s.Type {
	case GitHub:
		return GitHubBaseURL + s.Source + "/tree/" + ref
	default:
		return ""
	}
}

func (s *Repository) ActionRunURL(runID string) string {
	if runID == "" {
		return ""
	}

	switch s.Type {
	case GitHub:
		return GitHubBaseURL + s.Source + "/actions/runs/" + runID
	default:
		return ""
	}
}
