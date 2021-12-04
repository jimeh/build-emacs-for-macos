package plan

import (
	"context"
	"fmt"
	"io"
	"regexp"
	"strconv"
	"strings"

	"github.com/hashicorp/go-hclog"
	"github.com/jimeh/build-emacs-for-macos/pkg/commit"
	"github.com/jimeh/build-emacs-for-macos/pkg/gh"
	"github.com/jimeh/build-emacs-for-macos/pkg/osinfo"
	"github.com/jimeh/build-emacs-for-macos/pkg/release"
	"github.com/jimeh/build-emacs-for-macos/pkg/repository"
	"github.com/jimeh/build-emacs-for-macos/pkg/sanitize"
	"github.com/jimeh/build-emacs-for-macos/pkg/source"
)

var gitTagMatcher = regexp.MustCompile(
	`^emacs(-.*)?-((\d+\.\d+)(?:\.(\d+))?(-rc\d+)?(.+)?)$`,
)

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

func Create(ctx context.Context, opts *Options) (*Plan, error) { //nolint:funlen
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

	absoluteVersion := fmt.Sprintf(
		"%s.%s.%s",
		commitInfo.DateString(),
		commitInfo.ShortSHA(),
		sanitize.String(opts.Ref),
	)

	version, channel, err := parseGitRef(opts.Ref)
	if err != nil {
		return nil, err
	}

	var releaseName string
	switch channel {
	case release.Stable, release.RC:
		releaseName = "Emacs-" + version
	case release.Pretest:
		version += "-pretest"
		absoluteVersion += "-pretest"
		releaseName = "Emacs-" + version
	default:
		version = absoluteVersion
		releaseName = "Emacs." + version
	}

	buildName := fmt.Sprintf(
		"Emacs.%s.%s.%s",
		absoluteVersion,
		sanitize.String(osInfo.Name+"-"+osInfo.DistinctVersion()),
		sanitize.String(osInfo.Arch),
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
			Prerelease: channel != release.Stable,
			Channel:    channel,
		},
		Output: &Output{
			Directory: opts.OutputDir,
			DiskImage: diskImage,
		},
	}

	if opts.TestBuild != "" {
		testName := sanitize.String(opts.TestBuild)

		plan.Build.Name += ".test." + testName
		plan.Release.Title = "Test Builds (" + testName + ")"
		plan.Release.Name = "test-builds"

		plan.Release.Prerelease = false
		plan.Release.Draft = true
		if opts.TestBuildType == Prerelease {
			plan.Release.Prerelease = true
			plan.Release.Draft = false
		}

		index := strings.LastIndex(diskImage, ".")
		plan.Output.DiskImage = diskImage[:index] + ".test." +
			testName + diskImage[index:]
	}

	return plan, nil
}

func parseGitRef(ref string) (string, release.Channel, error) {
	m := gitTagMatcher.FindStringSubmatch(ref)

	if len(m) == 0 {
		return "", release.Nightly, nil
	}

	if strings.Contains(m[1], "pretest") {
		return m[2], release.Pretest, nil
	}

	if m[4] != "" {
		n, err := strconv.Atoi(m[4])
		if err != nil {
			return "", "", err
		}

		if n >= 90 {
			return m[2], release.Pretest, nil
		}
	}

	if strings.HasPrefix(m[5], "-rc") {
		return m[2], release.RC, nil
	}

	if m[2] == m[3] {
		return m[2], release.Stable, nil
	}

	return "", "", nil
}
