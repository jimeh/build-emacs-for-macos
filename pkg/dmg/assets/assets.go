package assets

import (
	_ "embed"
	"os"
)

//go:generate tiffutil -cathidpicheck bg.png bg@2x.png -out bg.tif

// Background is a raw byte slice of bytes of bg.tiff
//
//go:embed bg.tif
var Background []byte

// BackgroundTempFile writes Background to a temporary file on disk, returning
// the resulting file path. The returned filepath should be deleted with
// os.Remove() when no longer needed.
func BackgroundTempFile() (string, error) {
	return tempFile("*-emacs-bg.tif", Background)
}

// Icon is a raw byte slice of bytes of vol.icns
//
//go:embed vol.icns
var Icon []byte

// IconTempFile writes Icon to a temporary file on disk, returning the resulting
// file path. The returned filepath should be deleted with os.Remove() when no
// longer needed.
func IconTempFile() (string, error) {
	return tempFile("*-emacs-vol.icns", Icon)
}

func tempFile(pattern string, content []byte) (string, error) {
	f, err := os.CreateTemp("", pattern)
	if err != nil {
		return "", err
	}
	defer f.Close()

	_, err = f.Write(content)
	if err != nil {
		return "", err
	}

	return f.Name(), nil
}
