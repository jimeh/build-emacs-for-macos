package release

import (
	"bytes"
	"os"
	"strings"
	"text/template"
)

var tplFuncs = template.FuncMap{
	"indent": func(n int, s string) string {
		pad := strings.Repeat(" ", n)

		return pad + strings.ReplaceAll(s, "\n", "\n"+pad)
	},
}

var bodyTpl = template.Must(template.New("body").Funcs(tplFuncs).Parse(`
{{- $t := "` + "`" + `" -}}
### Build Details

{{ with .SourceURL -}}
- Source: {{ . }}
{{- end }}
{{- if .CommitURL }}
- Commit: {{ .CommitURL }}
  {{- if .CommitSHA }} ({{ $t }}{{ .CommitSHA }}{{ $t }}){{ end }}
{{- end }}
{{- with .TarballURL }}
- Tarball: {{ . }}
{{- end }}
{{- with .BuildLogURL }}
- Build Log: {{ . }} (available for 90 days)
{{- end }}`,
))

type bodyData struct {
	SourceURL   string
	CommitSHA   string
	CommitURL   string
	BuildLogURL string
	TarballURL  string
}

func releaseBody(opts *PublishOptions) (string, error) {
	src := opts.Source

	if src.Repository == nil || src.Commit == nil {
		return "", nil
	}

	data := &bodyData{
		SourceURL:  src.Repository.TreeURL(src.Ref),
		CommitSHA:  src.Commit.SHA,
		CommitURL:  src.Repository.CommitURL(src.Commit.SHA),
		TarballURL: src.Repository.TarballURL(src.Commit.SHA),
	}

	// If available, use the exact value from the build plan.
	if src.Tarball != nil {
		data.TarballURL = src.Tarball.URL
	}

	// If running within GitHub Actions, provide link to build log.
	if opts.Repository != nil {
		if id := os.Getenv("GITHUB_RUN_ID"); id != "" {
			data.BuildLogURL = opts.Repository.ActionRunURL(id)
		}
	}

	var buf bytes.Buffer
	err := bodyTpl.Execute(&buf, data)
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}
