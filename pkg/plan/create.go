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
	"github.com/jimeh/build-emacs-for-macos/pkg/repository"
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

	releaseName := fmt.Sprintf(
		"Emacs.%s.%s.%s",
		commitInfo.DateString(),
		commitInfo.ShortSHA(),
		sanitizeString(opts.Ref),
	)
	buildName := fmt.Sprintf(
		"%s.%s.%s",
		releaseName,
		sanitizeString(osInfo.Name+"-"+osInfo.MajorMinor()),
		sanitizeString(osInfo.Arch),
	)
	diskImage := buildName + ".dmg"

	plan := &Plan{
		Build: &Build{
			Name: buildName,
		},
		Source: &Source{
			Ref:        opts.Ref,
			Repository: repo,
			Commit:     commitInfo,
			Tarball: &Tarball{
				URL: repo.TarballURL(commitInfo.SHA),
			},
		},
		OS: osInfo,
		Release: &Release{
			Name: releaseName,
		},
		Output: &Output{
			Directory: opts.OutputDir,
			DiskImage: diskImage,
		},
	}

	if opts.TestBuild != "" {
		testName := sanitizeString(opts.TestBuild)

		plan.Build.Name += ".test." + testName
		plan.Release.Title = "Test Builds"
		plan.Release.Name = "test-builds"
		if opts.TestBuildType == Draft {
			plan.Release.Draft = true
		} else {
			plan.Release.Prerelease = true
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
