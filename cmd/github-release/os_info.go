package main

import (
	"os/exec"
	"strings"
)

var OSInfo = &osInfo{}

type osInfo struct {
	version string
	arch    string
}

func (s *osInfo) Arch() string {
	if s.arch != "" {
		return s.arch
	}

	cmd, err := exec.Command("uname", "-m").CombinedOutput()
	if err != nil {
		panic(err)
	}

	s.arch = strings.TrimSpace(string(cmd))

	return s.arch
}

func (s *osInfo) Version() string {
	if s.version != "" {
		return s.version
	}

	cmd, err := exec.Command("sw_vers", "-productVersion").CombinedOutput()
	if err != nil {
		panic(err)
	}

	parts := strings.Split(string(cmd), ".")
	s.version = strings.Join(parts[0:2], ".")

	return s.version
}
