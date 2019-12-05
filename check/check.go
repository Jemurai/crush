package check

// Check is a security check.
type Check struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Magic       string   `json:"magic"`
	Extensions  []string `json:"exts"`
	Tags        []string `json:"tags"`
	Ran         bool     `json:"ran"`
}

// CountHowManyRan - given an array of checks, how many ran.
func CountHowManyRan(checks []Check) int {
	ran := 0
	for i := 0; i < len(checks); i++ {
		if checks[i].Ran == true {
			ran = ran + 1
		}
	}
	return ran
}

// AppliesToTag checks that the specific check should be
// applied to a given tag.
func AppliesToTag(check Check, tag string) bool {
	if tag == "" {
		return true
	}
	for i := 0; i < len(check.Tags); i++ {
		if check.Tags[i] == tag {
			return true
		}
	}
	return false
}

// AppliesToExt checks that the extension provided should
// be checked.  (Eg. .java)
func AppliesToExt(check Check, ext string) bool {
	if ext == "" {
		return true
	}
	for i := 0; i < len(check.Extensions); i++ {
		if check.Extensions[i] == ext {
			return true
		}
	}
	return false
}
