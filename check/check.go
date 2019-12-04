package check

// Check is a security check.
type Check struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Magic       string   `json:"magic"`
	Extensions  []string `json:"exts"`
	Tags        []string `json:"tags"`
}
