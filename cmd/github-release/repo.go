package main

import "strings"

type Repo struct {
	Owner string
	Name  string
}

func NewRepo(ownerAndRepo string) *Repo {
	parts := strings.SplitN(ownerAndRepo, "/", 2)

	return &Repo{
		Owner: parts[0],
		Name:  parts[1],
	}
}

func (s *Repo) String() string {
	return s.Owner + "/" + s.Name
}
