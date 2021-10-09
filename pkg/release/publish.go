package release

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/google/go-github/v35/github"
	"github.com/hashicorp/go-hclog"
	"github.com/jimeh/build-emacs-for-macos/pkg/gh"
	"github.com/jimeh/build-emacs-for-macos/pkg/repository"
	"github.com/jimeh/build-emacs-for-macos/pkg/source"
)

type releaseType int

// Release type constants
const (
	Normal releaseType = iota
	Draft
	Prerelease
)

type PublishOptions struct {
	// Repository is the GitHub repository to publish the release on.
	Repository *repository.Repository

	// CommitRef is the git commit ref to create the release tag on.
	CommitRef string

	// ReleaseName is the name of the git tag to create the release on.
	ReleaseName string

	// ReleaseTitle is the title of the release, if not specified ReleaseName
	// will be used.
	ReleaseTitle string

	// ReleaseType is the type of release to create (normal, prerelease, or
	// draft)
	ReleaseType releaseType

	// Source contains the source used to build the asset files. When set a
	// release body/description text will be generated based on source commit
	// details.
	Source *source.Source

	// AssetFiles is a list of files which must all exist in the release for
	// the check to pass.
	AssetFiles []string

	// AssetSizeCheck causes a file size check for any existing asset files on a
	// release which have the same filename as a asset we want to upload. If the
	// size of the local and remote files are the same, the existing asset file
	// is left in place. When this is false, given asset files will always be
	// uploaded, replacing any asset files with the same filename.
	AssetSizeCheck bool

	// GitHubToken is the OAuth token used to talk to the GitHub API.
	GithubToken string
}

// Publish creates and publishes a GitHub release.
func Publish(ctx context.Context, opts *PublishOptions) error {
	logger := hclog.FromContext(ctx).Named("release")
	gh := gh.New(ctx, opts.GithubToken)

	files, err := publishFileList(opts.AssetFiles)
	if err != nil {
		return err
	}

	tagName := opts.ReleaseName
	name := opts.ReleaseTitle
	if name == "" {
		name = tagName
	}

	prerelease := opts.ReleaseType == Prerelease
	draft := opts.ReleaseType == Draft

	body := ""
	if opts.Source != nil {
		body, err = releaseBody(opts)
		if err != nil {
			return err
		}
		logger.Debug("rendered release body", "content", body)
	}

	created := false
	logger.Info("checking release", "tag", tagName)
	release, resp, err := gh.Repositories.GetReleaseByTag(
		ctx, opts.Repository.Owner(), opts.Repository.Name(), tagName,
	)
	if err != nil {
		if resp.StatusCode != http.StatusNotFound {
			return err
		}

		created = true
		logger.Info("creating release", "tag", tagName, "name", name)

		release, _, err = gh.Repositories.CreateRelease(
			ctx, opts.Repository.Owner(), opts.Repository.Name(),
			&github.RepositoryRelease{
				Name:            &name,
				TagName:         &tagName,
				TargetCommitish: &opts.CommitRef,
				Prerelease:      boolPtr(false),
				Draft:           boolPtr(true),
				Body:            &body,
			},
		)
		if err != nil {
			return err
		}
	}

	err = uploadReleaseAssets(ctx, gh, release, files, opts)
	if err != nil {
		return err
	}

	changed := false
	if release.GetName() != name {
		release.Name = &name
		changed = true
	}

	if body != "" && release.GetBody() != body {
		release.Body = &body
		changed = true
	}

	if release.GetDraft() != draft {
		release.Draft = &draft
		changed = true
	}

	if !draft && release.GetPrerelease() != prerelease {
		release.Prerelease = &prerelease
		changed = true
	}

	if changed {
		logger.Info("updating release attributes", "url", release.GetHTMLURL())
		release, _, err = gh.Repositories.EditRelease(
			ctx, opts.Repository.Owner(), opts.Repository.Name(),
			release.GetID(), release,
		)
		if err != nil {
			return err
		}
	}

	if created {
		logger.Info("release created", "url", release.GetHTMLURL())
	} else {
		logger.Info("release updated", "url", release.GetHTMLURL())
	}

	return nil
}

func uploadReleaseAssets(
	ctx context.Context,
	gh *github.Client,
	release *github.RepositoryRelease,
	fileNames []string,
	opts *PublishOptions,
) error {
	logger := hclog.FromContext(ctx).Named("release")

	for _, fileName := range fileNames {
		logger.Debug("processing asset", "file", filepath.Base(fileName))

		fileIO, err := os.Open(fileName)
		if err != nil {
			return err
		}
		defer fileIO.Close()

		fileInfo, err := fileIO.Stat()
		if err != nil {
			return err
		}

		fileBaseName := filepath.Base(fileName)
		assetExists := false

		for _, a := range release.Assets {
			if a.GetName() != fileBaseName {
				continue
			}

			if opts.AssetSizeCheck && a.GetSize() == int(fileInfo.Size()) {
				logger.Info("asset exists with correct size",
					"file", fileBaseName,
					"local_size", byteCountIEC(fileInfo.Size()),
					"remote_size", byteCountIEC(int64(a.GetSize())),
				)
				assetExists = true
			} else {
				logger.Info(
					"deleting existing asset", "file", fileBaseName,
				)
				_, err = gh.Repositories.DeleteReleaseAsset(
					ctx, opts.Repository.Owner(), opts.Repository.Name(),
					a.GetID(),
				)
				if err != nil {
					return err
				}
			}
		}

		if !assetExists {
			logger.Info("uploading asset",
				"file", fileBaseName,
				"size", byteCountIEC(fileInfo.Size()),
			)
			_, _, err = gh.Repositories.UploadReleaseAsset(
				ctx, opts.Repository.Owner(), opts.Repository.Name(),
				release.GetID(),
				&github.UploadOptions{Name: fileBaseName},
				fileIO,
			)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func publishFileList(files []string) ([]string, error) {
	var output []string
	for _, file := range files {
		var err error
		file, err = filepath.Abs(file)
		if err != nil {
			return nil, err
		}

		stat, err := os.Stat(file)
		if err != nil {
			return nil, err
		}
		if !stat.Mode().IsRegular() {
			return nil, fmt.Errorf("\"%s\" is not a file", file)
		}

		output = append(output, file)
		sumFile := file + ".sha256"

		_, err = os.Stat(sumFile)
		if err != nil {
			if os.IsNotExist(err) {
				continue
			}

			return nil, err
		}
		output = append(output, sumFile)
	}

	return output, nil
}

func byteCountIEC(b int64) string {
	const unit = 1024
	if b < unit {
		if b == 1 {
			return fmt.Sprintf("%d byte", b)
		}

		return fmt.Sprintf("%d bytes", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}

	return fmt.Sprintf("%.1f %ciB",
		float64(b)/float64(div), "KMGTPE"[exp])
}

func boolPtr(v bool) *bool {
	return &v
}
