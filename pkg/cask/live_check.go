package cask

type LiveCheck struct {
	Cask    string           `json:"cask"`
	Version LiveCheckVersion `json:"version"`
}

type LiveCheckVersion struct {
	Current           string `json:"current"`
	Latest            string `json:"latest"`
	Outdated          bool   `json:"outdated"`
	NewerThanUpstream bool   `json:"newer_than_upstream"`
}
