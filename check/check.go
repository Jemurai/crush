// Copyright © 2019-2023 Matt Konda <mkonda@jemurai.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package check

import (
	"encoding/json"
	"regexp"
	"strings"
)

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

func (c Check) String() string {
	s, _ := json.Marshal(c)
	return string(s)
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
	if len(extensions) == 0 {
		return true
	}
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
	}
	return false
}

// AppliesBasedOnComment tries to ignore common comment formats
// to avoid false positives where issues are in comments.
func AppliesBasedOnComment(line string, ext string) bool {
	trimmed := strings.TrimSpace(line)
	if checkMatch(trimmed, "crushignore") {
		return false
	} else if ext == ".rb" &&
		checkMatch(trimmed, "^#") {
		return false
	} else if ext == ".js" &&
		checkMatch(trimmed, "^//") {
		return false
	} else if ext == ".clj" &&
		checkMatch(trimmed, "^;") {
		return false
	} else if ext == ".ex" &&
		checkMatch(trimmed, "^#") {
		return false
	}
	return true
}

func checkMatch(line string, comment string) bool {
	r, _ := regexp.Compile(comment)
	matched := r.MatchString(line)
	return matched
}
