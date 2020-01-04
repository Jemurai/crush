package check

// Check is a security check.
type Check struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Magic       string   `json:"magic"`
	Extensions  []string `json:"exts"`
	Tags        []string `json:"tags"`
	Ran         bool     `json:"ran"`
	Threshold   float64  `json:"threshold"`
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

// AppliesToExt checks that the extension provided should be checked.
// For example does it apply to .java.
//
// There are two dimensions for this:
// 1.  Does the check apply to the extension of the actual file: actualExt
// 2.  Did the user specify they wanted to only run checks on certain extensions
func AppliesToExt(check Check, actualExt string, extOption string) bool {
	var applies bool
	if extOption == "" {
		applies = checkExtensions(check.Extensions, actualExt)
	} else {
		applies = checkExtensions(check.Extensions, extOption) && checkExtensions(check.Extensions, actualExt)
	}
	return applies
}

func checkExtensions(extensions []string, extension string) bool {
	for i := 0; i < len(extensions); i++ {
		if extensions[i] == extension {
			return true
		}
	}
	return false
}

// AppliesBasedOnThreshold checks that the check applies based on confidence
// selection
func AppliesBasedOnThreshold(check Check, threshold float64) bool {
	if check.Threshold >= threshold {
		return true
	} else {
		return false
	}
}
