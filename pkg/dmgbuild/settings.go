package dmgbuild

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/hashicorp/go-hclog"
)

type format string

//nolint:golint
var (
	UDROFormat format = "UDRO" // Read-only
	UDCOFormat format = "UDCO" // Compressed (ADC)
	UDZOFormat format = "UDZO" // Compressed (gzip)
	UDBZFormat format = "UDBZ" // Compressed (bzip2)
	UFBIFormat format = "UFBI" // Entire device
	IPODFormat format = "IPOD" // iPod image
	UDxxFormat format = "UDxx" // UDIF stub
	UDSBFormat format = "UDSB" // Sparse bundle
	UDSPFormat format = "UDSP" // Sparse
	UDRWFormat format = "UDRW" // Read/write
	UDTOFormat format = "UDTO" // DVD/CD master
	DC42Format format = "DC42" // Disk Copy 4.2
	RdWrFormat format = "RdWr" // NDIF read/write
	RdxxFormat format = "Rdxx" // NDIF read-only
	ROCoFormat format = "ROCo" // NDIF Compressed
	RkenFormat format = "Rken" // NDIF Compressed (KenCode)
)

type File struct {
	Path          string
	PosX          int
	PosY          int
	Hidden        bool
	HideExtension bool
}

type Symlink struct {
	Name          string
	Target        string
	PosX          int
	PosY          int
	Hidden        bool
	HideExtension bool
}

type Settings struct {
	// Command can be set to a custom dmgbuild executable path. If not set,
	// the first "dmgbuild" executable within PATH will be used.
	Command string

	// Stdout will be set as STDOUT target for dmgbuild execution if not nil.
	Stdout io.Writer

	// Stderr will be set as STDERR target for dmgbuild execution if not nil.
	Stderr io.Writer

	// Logger allows logging details of dmbuild process.
	Logger hclog.Logger

	// dmgbuild settings
	Filename         string
	VolumeName       string
	Format           format
	Size             string
	CompressionLevel int
	Files            []*File
	Symlinks         []*Symlink
	Icon             string
	BadgeIcon        string
	Window           Window
	IconView         IconView
	ListView         ListView
	License          License
}

func NewSettings() *Settings {
	return &Settings{
		Format:           UDZOFormat,
		CompressionLevel: 9,
		Window:           NewWindow(),
		IconView:         NewIconView(),
		ListView:         NewListView(),
		License:          NewLicense(),
	}
}

//nolint:funlen,gocyclo
// Render returns a string slice where each string is a separate settings
// statement.
func (s *Settings) Render() ([]string, error) {
	r := []string{
		"# -*- coding: utf-8 -*-\n",
		"from __future__ import unicode_literals\n",
	}

	if s.Filename != "" {
		r = append(r, "filename = "+pyStr(s.Filename)+"\n")
	}
	if s.VolumeName != "" {
		r = append(r, "volume_name = "+pyStr(s.VolumeName)+"\n")
	}
	if s.Format != "" {
		r = append(r, "format = "+pyStr(string(s.Format))+"\n")
	}
	if s.CompressionLevel != 0 {
		r = append(r, fmt.Sprintf(
			"compression_level = %d\n", s.CompressionLevel,
		))
	}
	if s.Size != "" {
		r = append(r, "size = "+pyStr(s.Size)+"\n")
	}

	var files []string
	var symlinks []string
	var hide []string
	var hideExt []string
	var iconLoc []string

	if len(s.Files) > 0 {
		for _, f := range s.Files {
			files = append(files, pyStr(f.Path))
			name := filepath.Base(f.Path)
			if f.PosX > 0 || f.PosY > 0 {
				iconLoc = append(iconLoc,
					fmt.Sprintf("%s: (%d, %d)", pyStr(name), f.PosX, f.PosY),
				)
			}
			if f.Hidden {
				hide = append(hide, pyStr(filepath.Base(f.Path)))
			}
			if f.HideExtension {
				hideExt = append(hideExt, pyStr(filepath.Base(f.Path)))
			}
		}
	}

	if len(s.Symlinks) > 0 {
		for _, l := range s.Symlinks {
			symlinks = append(symlinks, pyStr(l.Name)+": "+pyStr(l.Target))
			if l.PosX > 0 || l.PosY > 0 {
				iconLoc = append(iconLoc,
					fmt.Sprintf("%s: (%d, %d)", pyStr(l.Name), l.PosX, l.PosY),
				)
			}
			if l.Hidden {
				hide = append(hide, pyStr(l.Name))
			}
			if l.HideExtension {
				hideExt = append(hideExt, pyStr(l.Name))
			}
		}
	}

	if len(files) > 0 {
		r = append(r,
			"files = [\n    "+strings.Join(files, ",\n    ")+"\n]\n",
		)
	}
	if len(symlinks) > 0 {
		r = append(r,
			"symlinks = {\n    "+strings.Join(symlinks, ",\n    ")+"\n}\n",
		)
	}
	if len(hide) > 0 {
		r = append(r,
			"hide = [\n    "+strings.Join(hide, ",\n    ")+"\n]\n",
		)
	}
	if len(hideExt) > 0 {
		r = append(r,
			"hide_extensions = [\n    "+strings.Join(hideExt, ",\n    ")+
				"\n]\n",
		)
	}
	if len(iconLoc) > 0 {
		r = append(r,
			"icon_locations = {\n    "+strings.Join(iconLoc, ",\n    ")+"\n}\n",
		)
	}

	if s.Icon != "" {
		r = append(r, "icon = "+pyStr(s.Icon)+"\n")
	}
	if s.BadgeIcon != "" {
		r = append(r, "badge_icon = "+pyStr(s.BadgeIcon)+"\n")
	}

	r = append(r, s.Window.Render()...)
	r = append(r, s.IconView.Render()...)
	r = append(r, s.ListView.Render()...)
	r = append(r, s.License.Render()...)

	return r, nil
}

func (s *Settings) Write(w io.Writer) error {
	out, err := s.Render()
	if err != nil {
		return err
	}

	for _, o := range out {
		_, err := w.Write([]byte(o))
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *Settings) TempFile() (string, error) {
	f, err := os.CreateTemp("", "*.dmgbuild.settings.py")
	if err != nil {
		return "", err
	}
	defer f.Close()

	err = s.Write(f)
	if err != nil {
		return "", err
	}

	return f.Name(), nil
}

func pyStr(s string) string {
	s = strings.ReplaceAll(s, `\`, `\\`)
	s = strings.ReplaceAll(s, `"`, `\"`)
	s = strings.ReplaceAll(s, "\r", `\r`)
	s = strings.ReplaceAll(s, "\n", `\n`)

	return `"` + s + `"`
}

func pyMStr(s string) string {
	s = strings.ReplaceAll(s, `\`, `\\`)
	s = strings.ReplaceAll(s, `"`, `\"`)

	return `"""` + s + `"""`
}

func pyBool(v bool) string {
	if v {
		return "True"
	}

	return "False"
}
