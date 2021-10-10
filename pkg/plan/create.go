package plan

import (
	"context"
	"fmt"
	"io"
	"regexp"
	"strings"

	"github.com/hashicorp/go-hclog"
	"github.com/jimeh/build-emacs-for-macos/pkg/commit"
	"github.com/jimeh/build-emacs-for-macos/pkg/gh"
	"github.com/jimeh/build-emacs-for-macos/pkg/osinfo"
	"github.com/jimeh/build-emacs-for-macos/pkg/release"
	"github.com/jimeh/build-emacs-for-macos/pkg/repository"
	"github.com/jimeh/build-emacs-for-macos/pkg/source"
)

var nonAlphaNum = regexp.MustCompile(`[^\w_-]+`)

type TestBuildType string

//nolint:golint
const (
	Draft      TestBuildType = "draft"
	Prerelease TestBuildType = "prerelease"
)

type Options struct {
	GithubToken   string
	EmacsRepo     string
	Ref           string
	SHAOverride   string
	OutputDir     string
	TestBuild     string
	TestBuildType TestBuildType
	Output        io.Writer
}

func Create(ctx context.Context, opts *Options) (*Plan, error) {
	logger := hclog.FromContext(ctx).Named("plan")

	repo, err := repository.NewGitHub(opts.EmacsRepo)
	if err != nil {
		return nil, err
	}

	gh := gh.New(ctx, opts.GithubToken)

	lookupRef := opts.Ref
	if opts.SHAOverride != "" {
		lookupRef = opts.SHAOverride
	}
	logger.Info("fetching commit info", "ref", lookupRef)

	repoCommit, _, err := gh.Repositories.GetCommit(
		ctx, repo.Owner(), repo.Name(), lookupRef,
	)
	if err != nil {
		return nil, err
	}

	commitInfo := commit.New(repoCommit)
	osInfo, err := osinfo.New()
	if err != nil {
		return nil, err
	}

	version := fmt.Sprintf(
		"%s.%s.%s",
		commitInfo.DateString(),
		commitInfo.ShortSHA(),
		sanitizeString(opts.Ref),
	)

	releaseName := fmt.Sprintf("Emacs.%s", version)
	buildName := fmt.Sprintf(
		"Emacs.%s.%s.%s",
		version,
		sanitizeString(osInfo.Name+"-"+osInfo.DistinctVersion()),
		sanitizeString(osInfo.Arch),
	)
	diskImage := buildName + ".dmg"

	plan := &Plan{
		Build: &Build{
			Name: buildName,
		},
		Source: &source.Source{
			Ref:        opts.Ref,
			Repository: repo,
			Commit:     commitInfo,
			Tarball: &source.Tarball{
				URL: repo.TarballURL(commitInfo.SHA),
			},
		},
		OS: osInfo,
		Release: &Release{
			Name:       releaseName,
			Prerelease: true,
		},
		Output: &Output{
			Directory: opts.OutputDir,
			DiskImage: diskImage,
		},
	}

	// If given git ref is a stable release tag (emacs-23.2b, emacs-27.2, etc.)
	// we modify release properties accordingly.
	if v, err := release.GitRefToStableVersion(opts.Ref); err == nil {
		plan.Release.Prerelease = false
		plan.Release.Name, err = release.VersionToName(v)
		if err != nil {
			return nil, err
		}
	}

	if opts.TestBuild != "" {
		testName := sanitizeString(opts.TestBuild)

		plan.Build.Name += ".test." + testName
		plan.Release.Title = "Test Builds"
		plan.Release.Name = "test-builds"

		plan.Release.Prerelease = true
		plan.Release.Draft = false
		if opts.TestBuildType == Draft {
			plan.Release.Prerelease = false
			plan.Release.Draft = true
		}

		index := strings.LastIndex(diskImage, ".")
		plan.Output.DiskImage = diskImage[:index] + ".test." +
			testName + diskImage[index:]
	}

	return plan, nil
}

func sanitizeString(s string) string {
	return nonAlphaNum.ReplaceAllString(s, "-")
}
