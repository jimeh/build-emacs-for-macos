package cask

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"time"

	"github.com/google/go-github/v35/github"
	"github.com/hashicorp/go-hclog"
	"github.com/hexops/gotextdiff"
	"github.com/hexops/gotextdiff/myers"
	"github.com/hexops/gotextdiff/span"
	"github.com/jimeh/build-emacs-for-macos/pkg/gh"
	"github.com/jimeh/build-emacs-for-macos/pkg/release"
	"github.com/jimeh/build-emacs-for-macos/pkg/repository"
)

// Error vars
var (
	Err                = errors.New("cask")
	ErrReleaseNotFound = fmt.Errorf("%w: release not found", Err)

	ErrFailedSHA256Parse = fmt.Errorf(
		"%w: failed to parse SHA256 from asset", Err,
	)
	ErrFailedSHA256Download = fmt.Errorf(
		"%w: failed to download SHA256 asset", Err,
	)
	ErrNoTapOrOutput = fmt.Errorf(
		"%w: no tap repository or output directory specified", Err,
	)
)

type UpdateOptions struct {
	// BuildsRepo is the GitHub repository containing binary releases.
	BuildsRepo *repository.Repository

	// TapRepo is the GitHub repository to update the casks in.
	TapRepo *repository.Repository

	// Ref is the git ref to apply cask updates on top of. Default branch will
	// be used if empty.
	Ref string

	// OutputDir specifies a directory to write cask files to. When set, tap
	// repository is ignored and no changes will be committed directly against
	// any specified tap repository.
	OutputDir string

	// Force update will ignore the outdated live check flag, and process all
	// casks regardless. But it will only update the cask in question if the
	// resulting output cask is different.
	Force bool

	// TemplatesDir is the directory where cask templates are located.
	TemplatesDir string

	LiveChecks []*LiveCheck

	GithubToken string
}

type Updater struct {
	BuildsRepo   *repository.Repository
	TapRepo      *repository.Repository
	Ref          string
	OutputDir    string
	TemplatesDir string

	logger hclog.Logger
	gh     *github.Client
}

func Update(ctx context.Context, opts *UpdateOptions) error {
	updater := &Updater{
		BuildsRepo:   opts.BuildsRepo,
		TapRepo:      opts.TapRepo,
		Ref:          opts.Ref,
		OutputDir:    opts.OutputDir,
		TemplatesDir: opts.TemplatesDir,
		logger:       hclog.FromContext(ctx).Named("cask"),
		gh:           gh.New(ctx, opts.GithubToken),
	}

	for _, chk := range opts.LiveChecks {
		err := updater.Update(ctx, chk, opts.Force)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *Updater) Update(
	ctx context.Context,
	chk *LiveCheck,
	force bool,
) error {
	if s.TapRepo == nil && s.OutputDir == "" {
		return ErrNoTapOrOutput
	}

	if !force && !chk.Version.Outdated {
		s.logger.Info("skipping", "cask", chk.Cask, "reason", "up to date")

		return nil
	}

	newCaskContent, err := s.renderCask(ctx, chk)
	if err != nil {
		return err
	}

	caskFile := chk.Cask + ".rb"

	if s.OutputDir != "" {
		_, err = s.putFile(
			ctx, chk, filepath.Join(s.OutputDir, caskFile), newCaskContent,
		)
		if err != nil {
			return err
		}

		return nil
	}

	_, err = s.putRepoFile(
		ctx, s.TapRepo, s.Ref, chk,
		filepath.Join("Casks", caskFile), newCaskContent,
	)
	if err != nil {
		return err
	}

	return nil
}

func (s *Updater) putFile(
	_ context.Context,
	chk *LiveCheck,
	filename string,
	content []byte,
) (bool, error) {
	parent := filepath.Dir(filename)
	s.logger.Info("processing cask update",
		"output-directory", parent, "cask", chk.Cask, "file", filename,
	)

	err := os.MkdirAll(parent, 0o755)
	if err != nil {
		return false, err
	}

	existingContent, err := os.ReadFile(filename)
	if err != nil && !os.IsNotExist(err) {
		return false, err
	}

	infoMsg := "creating cask"

	if !os.IsNotExist(err) {
		infoMsg = "updating cask"
		if bytes.Equal(existingContent, content) {
			s.logger.Info(
				"skip update: no change to cask content",
				"cask", chk.Cask, "file", filename,
			)

			s.logger.Debug(
				"cask content",
				"file", filename, "content", string(content),
			)

			return false, nil
		}
	}

	existing := string(existingContent)
	edits := myers.ComputeEdits(
		span.URIFromPath(filename), existing, string(content),
	)
	diff := fmt.Sprint(gotextdiff.ToUnified(
		filename, filename, existing, edits,
	))

	s.logger.Info(
		infoMsg,
		"cask", chk.Cask, "version", chk.Version.Latest, "file", filename,
		"diff", diff,
	)

	s.logger.Debug(
		"cask content",
		"file", filename, "content", string(content),
	)

	err = os.WriteFile(filename, content, 0o644) //nolint:gosec
	if err != nil {
		return false, err
	}

	return true, nil
}

func (s *Updater) putRepoFile(
	ctx context.Context,
	repo *repository.Repository,
	ref string,
	chk *LiveCheck,
	filename string,
	content []byte,
) (bool, error) {
	s.logger.Info("processing cask update",
		"tap-repo", repo.Source, "cask", chk.Cask, "file", filename,
	)
	repoContent, _, resp, err := s.gh.Repositories.GetContents(
		ctx, repo.Owner(), repo.Name(), filename,
		&github.RepositoryContentGetOptions{Ref: ref},
	)
	if err != nil && resp.StatusCode != http.StatusNotFound {
		return false, err
	}

	if resp.StatusCode == http.StatusNotFound {
		err := s.createRepoFile(ctx, repo, chk, filename, content)
		if err != nil {
			return false, err
		}
	} else {
		_, err := s.updateRepoFile(
			ctx, repo, repoContent, chk, filename, content,
		)
		if err != nil {
			return false, err
		}
	}

	return true, nil
}

func (s *Updater) createRepoFile(
	ctx context.Context,
	repo *repository.Repository,
	chk *LiveCheck,
	filename string,
	content []byte,
) error {
	commitMsg := fmt.Sprintf(
		"feat(cask): create %s with version %s",
		chk.Cask, chk.Version.Latest,
	)

	edits := myers.ComputeEdits(
		span.URIFromPath(filename), "", string(content),
	)
	diff := fmt.Sprint(gotextdiff.ToUnified(filename, filename, "", edits))

	s.logger.Info(
		"creating cask",
		"cask", chk.Cask, "version", chk.Version.Latest, "file", filename,
		"diff", diff,
	)
	s.logger.Debug(
		"cask content",
		"file", filename, "content", string(content),
	)
	contResp, _, err := s.gh.Repositories.CreateFile(
		ctx, repo.Owner(), repo.Name(), filename,
		&github.RepositoryContentFileOptions{
			Message: &commitMsg,
			Content: content,
		},
	)
	if err != nil {
		return err
	}

	s.logger.Info(
		"new commit created",
		"commit", contResp.GetSHA(), "message", contResp.GetMessage(),
		"url", contResp.Commit.GetHTMLURL(),
	)

	return nil
}

func (s *Updater) updateRepoFile(
	ctx context.Context,
	repo *repository.Repository,
	repoContent *github.RepositoryContent,
	chk *LiveCheck,
	filename string,
	content []byte,
) (bool, error) {
	existingContent, err := repoContent.GetContent()
	if err != nil {
		return false, err
	}

	if existingContent == string(content) {
		s.logger.Info(
			"skip update: no change to cask content",
			"cask", chk.Cask, "file", filename,
		)

		return false, nil
	}

	sha := repoContent.GetSHA()

	commitMsg := fmt.Sprintf(
		"feat(cask): update %s to version %s",
		chk.Cask, chk.Version.Latest,
	)

	edits := myers.ComputeEdits(
		span.URIFromPath(filename), existingContent, string(content),
	)
	diff := fmt.Sprint(gotextdiff.ToUnified(
		filename, filename, existingContent, edits,
	))

	s.logger.Info(
		"updating cask",
		"cask", chk.Cask, "version", chk.Version.Latest, "file", filename,
		"diff", diff,
	)
	s.logger.Debug(
		"cask content",
		"file", filename, "content", string(content),
	)

	contResp, _, err := s.gh.Repositories.CreateFile(
		ctx, repo.Owner(), repo.Name(), filename,
		&github.RepositoryContentFileOptions{
			Message: &commitMsg,
			Content: content,
			SHA:     &sha,
		},
	)
	if err != nil {
		return false, err
	}

	s.logger.Info(
		"new commit created",
		"commit", contResp.GetSHA(), "message", contResp.GetMessage(),
		"url", contResp.Commit.GetHTMLURL(),
	)

	return true, nil
}

//nolint:funlen
func (s *Updater) renderCask(
	ctx context.Context,
	chk *LiveCheck,
) ([]byte, error) {
	releaseName, err := release.VersionToName(chk.Version.Latest)
	if err != nil {
		return nil, err
	}

	s.logger.Info("fetching release details",
		"release", releaseName, "repo", s.BuildsRepo.URL(),
	)
	release, resp, err := s.gh.Repositories.GetReleaseByTag(
		ctx, s.BuildsRepo.Owner(), s.BuildsRepo.Name(), releaseName,
	)
	if err != nil {
		return nil, err
	}
	if release == nil || resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("%w: %s", ErrReleaseNotFound, releaseName)
	}

	info := &ReleaseInfo{
		Name:    release.GetName(),
		Version: chk.Version.Latest,
		Assets:  map[string]*ReleaseAsset{},
	}

	s.logger.Info("processing release assets")
	for _, asset := range release.Assets {
		filename := asset.GetName()
		s.logger.Debug("processing asset", "filename", filename)

		filename = strings.TrimSuffix(filename, ".sha256")

		if _, ok := info.Assets[filename]; !ok {
			info.Assets[filename] = &ReleaseAsset{
				Filename: filename,
			}
		}

		if strings.HasSuffix(asset.GetName(), ".sha256") {
			s.logger.Debug("downloading *.sha256 asset to extract SHA256 value")
			r, err2 := s.downloadAssetContent(ctx, asset)
			if err2 != nil {
				return nil, err2
			}
			defer r.Close()

			content := make([]byte, 64)
			n, err2 := io.ReadAtLeast(r, content, 64)
			if err2 != nil {
				return nil, err2
			}
			if n < 64 {
				return nil, fmt.Errorf(
					"%w: %s", ErrFailedSHA256Parse, asset.GetName(),
				)
			}

			sha := string(content)[0:64]
			if sha == "" {
				return nil, fmt.Errorf(
					"%w: %s", ErrFailedSHA256Parse, asset.GetName(),
				)
			}

			info.Assets[filename].SHA256 = sha
		} else {
			info.Assets[filename].DownloadURL = asset.GetBrowserDownloadURL()
		}
	}

	tplContent, err := os.ReadFile(
		filepath.Join(s.TemplatesDir, chk.Cask+".rb.tpl"),
	)
	if err != nil {
		return nil, err
	}

	helperContent, err := os.ReadFile(
		filepath.Join(s.TemplatesDir, "_helpers.tpl"),
	)
	if err != nil && !os.IsNotExist(err) {
		return nil, err
	}

	if len(helperContent) > 0 {
		tplContent = append(helperContent, tplContent...)
	}

	tpl, err := template.New(chk.Cask).Parse(string(tplContent))
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	err = tpl.Execute(&buf, info)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (s *Updater) downloadAssetContent(
	ctx context.Context,
	asset *github.ReleaseAsset,
) (io.ReadCloser, error) {
	httpClient := &http.Client{Timeout: 60 * time.Second}

	r, downloadURL, err := s.gh.Repositories.DownloadReleaseAsset(
		ctx, s.BuildsRepo.Owner(), s.BuildsRepo.Name(),
		asset.GetID(), httpClient,
	)
	if err != nil {
		return nil, err
	}

	if r == nil && downloadURL != "" {
		req, err := http.NewRequestWithContext(ctx, "GET", downloadURL, nil)
		if err != nil {
			return nil, err
		}

		//nolint:bodyclose
		resp, err := httpClient.Do(req)
		if err != nil {
			return nil, err
		}
		r = resp.Body
	}

	if r == nil {
		return nil, fmt.Errorf(
			"%s: %s", ErrFailedSHA256Download, asset.GetName(),
		)
	}

	return r, nil
}
