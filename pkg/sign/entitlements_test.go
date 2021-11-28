package sign

import (
	"bytes"
	"os"
	"strings"
	"testing"

	"github.com/jimeh/undent"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var entitlementsTestCases = []struct {
	name         string
	entitlements Entitlements
	want         string
}{
	{
		name:         "none",
		entitlements: Entitlements{},
		//nolint:lll
		want: undent.String(`
            <?xml version="1.0" encoding="UTF-8"?>
            <!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
            <plist version="1.0">
              <dict>
              </dict>
            </plist>`,
		),
	},
	{
		name:         "one",
		entitlements: Entitlements{"com.apple.security.cs.allow-jit"},
		//nolint:lll
		want: undent.String(`
            <?xml version="1.0" encoding="UTF-8"?>
            <!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
            <plist version="1.0">
              <dict>
                <key>com.apple.security.cs.allow-jit</key>
                <true/>
              </dict>
            </plist>`,
		),
	},
	{
		name: "many",
		entitlements: Entitlements{
			"com.apple.developer.mail-client",
			"com.apple.developer.web-browser",
			"com.apple.security.automation.apple-events",
			"com.apple.security.cs.allow-dyld-environment-variables",
			"com.apple.security.cs.allow-jit",
			"com.apple.security.cs.disable-library-validation",
			"com.apple.security.network.client",
			"com.apple.security.network.server",
		},
		//nolint:lll
		want: undent.String(`
            <?xml version="1.0" encoding="UTF-8"?>
            <!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
            <plist version="1.0">
              <dict>
                <key>com.apple.developer.mail-client</key>
                <true/>
                <key>com.apple.developer.web-browser</key>
                <true/>
                <key>com.apple.security.automation.apple-events</key>
                <true/>
                <key>com.apple.security.cs.allow-dyld-environment-variables</key>
                <true/>
                <key>com.apple.security.cs.allow-jit</key>
                <true/>
                <key>com.apple.security.cs.disable-library-validation</key>
                <true/>
                <key>com.apple.security.network.client</key>
                <true/>
                <key>com.apple.security.network.server</key>
                <true/>
              </dict>
            </plist>`,
		),
	},
}

func TestDefaultEmacsEntitlements(t *testing.T) {
	assert.Equal(t,
		[]string{
			"com.apple.developer.mail-client",
			"com.apple.developer.web-browser",
			"com.apple.security.automation.apple-events",
			"com.apple.security.cs.allow-dyld-environment-variables",
			"com.apple.security.cs.allow-jit",
			"com.apple.security.cs.disable-library-validation",
			"com.apple.security.network.client",
			"com.apple.security.network.server",
		},
		DefaultEmacsEntitlements,
	)
}

func TestEntitlements_Write(t *testing.T) {
	for _, tt := range entitlementsTestCases {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer

			err := tt.entitlements.Write(&buf)
			require.NoError(t, err)

			assert.Equal(t, tt.want, strings.TrimSpace(buf.String()))
		})
	}
}

func TestEntitlements_TempFile(t *testing.T) {
	for _, tt := range entitlementsTestCases {
		t.Run(tt.name, func(t *testing.T) {
			tmpFile, err := tt.entitlements.TempFile()
			require.NoError(t, err)
			defer os.Remove(tmpFile)

			content, err := os.ReadFile(tmpFile)
			require.NoError(t, err)

			assert.Equal(t, tt.want, strings.TrimSpace(string(content)))
			assert.True(t,
				strings.HasSuffix(tmpFile, ".entitlements.plist"),
				"temp file name does not match \"*.entitlements.plist\"",
			)
		})
	}
}
