package sign

import "io"

type Options struct {
	Identity         string
	Entitlements     *Entitlements
	EntitlementsFile string
	Options          []string
	Deep             bool
	Timestamp        bool
	Force            bool
	Verbose          bool
	Output           io.Writer
	CodeSignCmd      string
}
