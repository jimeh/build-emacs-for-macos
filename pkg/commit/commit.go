package commit

import (
	"fmt"
	"time"

	"github.com/google/go-github/v35/github"
)

type Commit struct {
	SHA       string     `yaml:"sha"`
	Date      *time.Time `yaml:"date"`
	Author    string     `yaml:"author"`
	Committer string     `yaml:"committer"`
	Message   string     `yaml:"message"`
}

func New(rc *github.RepositoryCommit) *Commit {
	return &Commit{
		SHA:  rc.GetSHA(),
		Date: rc.GetCommit().GetCommitter().Date,
		Author: fmt.Sprintf(
			"%s <%s>",
			rc.GetCommit().GetAuthor().GetName(),
			rc.GetCommit().GetAuthor().GetEmail(),
		),
		Committer: fmt.Sprintf(
			"%s <%s>",
			rc.GetCommit().GetCommitter().GetName(),
			rc.GetCommit().GetCommitter().GetEmail(),
		),
		Message: rc.GetCommit().GetMessage(),
	}
}

func (s *Commit) ShortSHA() string {
	return s.SHA[0:7]
}

func (s *Commit) DateString() string {
	return s.Date.Format("2006-01-02")
}
