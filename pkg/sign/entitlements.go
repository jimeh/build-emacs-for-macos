package sign

import (
	"bytes"
	_ "embed"
	"io"
	"os"
	"text/template"
)

// DefaultEmacsEntitlements is the default set of entitlements application
// bundles are signed with if no entitlements are provided.
var DefaultEmacsEntitlements = []string{
	"com.apple.developer.mail-client",
	"com.apple.developer.web-browser",
	"com.apple.security.automation.apple-events",
	"com.apple.security.cs.allow-dyld-environment-variables",
	"com.apple.security.cs.allow-jit",
	"com.apple.security.cs.disable-library-validation",
	"com.apple.security.network.client",
	"com.apple.security.network.server",
}

//go:embed entitlements.tpl
var entitlementsTemplate string

type Entitlements []string

func (e Entitlements) XML() ([]byte, error) {
	var buf bytes.Buffer
	err := e.Write(&buf)

	return buf.Bytes(), err
}

func (e Entitlements) Write(w io.Writer) error {
	tpl, err := template.New("entitlements.plist").Parse(entitlementsTemplate)
	if err != nil {
		return err
	}

	return tpl.Execute(w, e)
}

func (e Entitlements) TempFile() (string, error) {
	f, err := os.CreateTemp("", "*.entitlements.plist")
	if err != nil {
		return "", err
	}
	defer f.Close()

	err = e.Write(f)
	if err != nil {
		return "", err
	}

	return f.Name(), nil
}
