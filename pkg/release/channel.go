package release

type Channel string

// Release channels
const (
	Stable  Channel = "stable"
	RC      Channel = "release-candidate"
	Pretest Channel = "pretest"
	Nightly Channel = "nightly"
)
