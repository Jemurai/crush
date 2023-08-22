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
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strconv"

	"embed"

	log "github.com/sirupsen/logrus"

	"github.com/jemurai/crush/options"
	"github.com/jemurai/crush/utils"
	"github.com/jemurai/fkit/finding"
)

//go:embed checks/*
var checkfs embed.FS

// GetChecks - get the checks from the specified file.
// Use either the checks embedded file system (see above)
// or the local file system in case someone wants
// to "bring their own" checks.
func GetChecks(file string) []Check {
	var checks []Check
	// This loads the packaged items
	data, err := checkfs.ReadFile(file)

	// This (in theory) loads the
	if err != nil || data == nil {
		log.Debug("Didn't find with packaged checks, looking for OS checks")
		log.Debug(err)
		_, err := os.Stat(file) // In regular file system, is at check/<name>.json.
		if os.IsNotExist(err) { // In Docker, is at:  /app/check/<name>.json
			file = "/app/" + file
		}
		log.Debugf("Getting checks for: %v", file)
		rfile, err := os.Open(file)
		if err != nil {
			log.Error(err)
		}
		bytes, err := ioutil.ReadAll(rfile)
		if err != nil {
			log.Error(err)
		}
		json.Unmarshal(bytes, &checks)
	} else {
		json.Unmarshal(data, &checks)
	}
	return checks
}

// GetFiles - get the files from the options specified.
func GetFiles(options options.Options) []string {
	var files []string
	err := filepath.Walk(options.Directory,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.IsDir() {
				files = append(files, path)
			}
			return nil
		})
	if err != nil {
		log.Error(err)
	}
	return files
}

// ReadLines - read lines out of a file
func ReadLines(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}

// DoCheck - do a particular check
func DoCheck(file string, lineno int, check Check, line string, options options.Options) []finding.Finding {
	var findings []finding.Finding

	// log.Debugf("Running check: %v on file %v", check.Name, file)

	check.Ran = true // TODO: Figure out how to track this.

	r, _ := regexp.Compile(check.Magic)
	if lineno == 1 { // Do this once per file - for things like file name.
		matched := r.MatchString(file)
		if matched {
			finding := finding.Finding{
				Name:        check.Name,
				Description: check.Description,
				Detail:      file,
				Source:      check.Name,
				Location:    file,
				Fingerprint: utils.Fingerprint(check.Name + ":" + file),
			}
			log.Debugf("Finding: %v", finding)
			findings = append(findings, finding)
		}
	}

	matched := r.MatchString(line)
	if matched {
		finding := finding.Finding{
			Name:        check.Name,
			Description: check.Description,
			Detail:      utils.Truncate(line, 320),
			Source:      check.Name,
			Location:    file + ":" + strconv.Itoa(lineno),
			Fingerprint: utils.Fingerprint(check.Name + file + strconv.Itoa(lineno) + line),
		}
		log.Debugf("Finding: %v", finding)
		findings = append(findings, finding)
	}

	return findings
}

// ProcessFile - process a file
func ProcessFile(fn string, checks []Check, options options.Options) []finding.Finding {
	log.Debugf("Processing %s", fn)
	lines, err := ReadLines(fn)
	if err != nil {
		log.Errorf("Error reading file %s", fn)
	}
	var findings []finding.Finding
	for i := 0; i < len(lines); i++ {
		for j := 0; j < len(checks); j++ {
			if AppliesToTag(checks[j], options.Tag) &&
				AppliesToExt(checks[j], filepath.Ext(fn), options.Ext) &&
				AppliesBasedOnThreshold(checks[j], options.Threshold) &&
				AppliesBasedOnComment(lines[i], filepath.Ext(fn)) {
				fs := DoCheck(fn, i, checks[j], lines[i], options)
				findings = append(findings, fs...)
			}
		}
	}
	// log.Debugf("\tProcessed %s", fn)
	return findings
}

// PrintFindings - Print out the findings
func PrintFindings(findings []finding.Finding, checks []Check, files []string) {
	if len(findings) > 0 {
		fjson, _ := json.MarshalIndent(findings, "", " ")
		fmt.Printf("%s", fjson)
		howMany := CountHowManyRan(checks)
		log.Debugf("\n\nSummary:\n\tFiles processed: %v \n\tChecks: %v (Not working yet) \n\tIssues: %v", len(files), howMany, len(findings))
	} else {
		fmt.Print("[]")
	}
}
