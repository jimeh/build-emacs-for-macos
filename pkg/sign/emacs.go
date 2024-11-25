package sign

import (
	"context"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/hashicorp/go-hclog"
)

// Emacs signs a Emacs.app application bundle with Apple's codesign utility,
// using correct default entitlements, and also pre-signing any *.eln files
// which are in the bundle, as codesign will not detect them as requiring
// signing even with the --deep flag.
func Emacs(ctx context.Context, appBundle string, opts *Options) error {
	if !strings.HasSuffix(appBundle, ".app") {
		return fmt.Errorf("%s is not a .app application bundle", appBundle)
	}

	appBundle, err := filepath.Abs(appBundle)
	if err != nil {
		return err
	}

	_, err = os.Stat(appBundle)
	if err != nil {
		return err
	}

	logger := hclog.FromContext(ctx).Named("sign")
	logger.Info("preparing to sign Emacs.app", "app", appBundle)

	newOpts := *opts

	if newOpts.EntitlementsFile == "" {
		if newOpts.Entitlements == nil {
			e := Entitlements(DefaultEmacsEntitlements)
			newOpts.Entitlements = &e
		}

		f, err2 := newOpts.Entitlements.TempFile()
		if err2 != nil {
			return err2
		}
		defer os.Remove(f)

		newOpts.EntitlementsFile = f
		newOpts.Entitlements = nil
	}

	err = signElnFiles(ctx, appBundle, &newOpts)
	if err != nil {
		return err
	}

	err = signCLIHelper(ctx, appBundle, &newOpts)
	if err != nil {
		return err
	}

	// Ensure app bundle is signed last, as modifications to the bundle after
	// signing will invalidate the signature. Hence anything within it that
	// needs to be separately signed, has to happen before signing the whole
	// application bundle.
	return Files(ctx, []string{appBundle}, &newOpts)
}

func signElnFiles(ctx context.Context, appBundle string, opts *Options) error {
	logger := hclog.FromContext(ctx).Named("sign")

	elnFiles, err := elnFiles(appBundle)
	if err != nil {
		return err
	}

	if len(elnFiles) == 0 {
		return nil
	}

	logger.Info(fmt.Sprintf(
		"found %d native-lisp *.eln files in %s to sign",
		len(elnFiles), filepath.Base(appBundle),
	))
	for _, file := range elnFiles {
		err := Files(ctx, []string{file}, opts)
		if err != nil {
			return err
		}
	}

	return nil
}

func signCLIHelper(ctx context.Context, appBundle string, opts *Options) error {
	logger := hclog.FromContext(ctx).Named("sign")

	cliHelper := filepath.Join(appBundle, "Contents", "MacOS", "bin", "emacs")
	fi, err := os.Stat(cliHelper)
	if err != nil && !os.IsNotExist(err) {
		return err
	} else if err == nil && fi.Mode().IsRegular() {
		logger.Info(fmt.Sprintf(
			"found Contents/MacOS/bin/emacs CLI helper script in %s to sign",
			filepath.Base(appBundle),
		))

		err = Files(ctx, []string{cliHelper}, opts)
		if err != nil {
			return err
		}
	}

	return nil
}

// elnFiles finds all native-compilation *.eln files within a Emacs.app bundle,
// excluding any *.eln which should be automatically located by codesign when
// signing the Emacs.app bundle itself with the --deep flag. Essentially this
// only returns *.eln files which must be individually signed before signing the
// app bundle itself.
func elnFiles(emacsApp string) ([]string, error) {
	var files []string
	walkDirFunc := func(path string, d fs.DirEntry, _ error) error {
		if d.Type().IsRegular() && strings.HasSuffix(path, ".eln") &&
			!strings.Contains(path, ".app/Contents/Frameworks/") {
			files = append(files, path)
		}

		return nil
	}

	err := filepath.WalkDir(filepath.Join(emacsApp, "Contents"), walkDirFunc)
	if err != nil {
		return nil, err
	}

	return files, nil
}
