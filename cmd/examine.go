// Copyright © 2019 Matt Konda <mkonda@jemurai.com>
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

package cmd

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"sync"
	"time"

	"github.com/jemurai/crush/check"
	"github.com/jemurai/crush/options"
	"github.com/jemurai/crush/utils"
	"github.com/jemurai/fkit/finding"

	log "github.com/sirupsen/logrus"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// examineCmd represents the share command
var examineCmd = &cobra.Command{
	Use:   "examine",
	Short: "Examine a directory",
	Long: `Examine all the source code in a directory.
	
Behind the scenes, crh checks a variety of things.`,

	Run: func(cmd *cobra.Command, args []string) {
		start := time.Now()
		opts := buildExamineOptions(cmd)
		print("Examine command on: " + opts.Directory)
		checks := getAllChecks(opts)
		files := getFiles(opts)
		var findings []finding.Finding
		var wg sync.WaitGroup
		for i := 0; i < len(files); i++ {
			wg.Add(1)
			fn := files[i]
			go func(fn string, opts options.Options) {
				defer wg.Done()
				findings = append(findings, processFile(fn, checks, opts)...)
			}(fn, opts)
		}
		wg.Wait()
		fjson, _ := json.MarshalIndent(findings, "", " ")
		fmt.Printf("[%s]", fjson)
		utils.Timing(start, "Elasped time: %f")
	},
}

func processFile(fn string, checks []check.Check, options options.Options) []finding.Finding {
	// log.Debugf("Processing %s", fn)
	lines, err := readLines(fn)
	if err != nil {
		log.Errorf("Error reading file %s", fn)
	}
	var findings []finding.Finding
	for i := 0; i < len(lines); i++ {
		for j := 0; j < len(checks); j++ {
			fs := doCheck(fn, i, checks[j], lines[i])
			findings = append(findings, fs...)
		}
	}
	// log.Debugf("\tProcessed %s", fn)
	return findings
}

func doCheck(file string, lineno int, check check.Check, line string) []finding.Finding {
	var findings []finding.Finding

	r, _ := regexp.Compile(check.Magic)
	if lineno == 1 { // Do this once per file.
		matched := r.MatchString(file)
		if matched {
			finding := finding.Finding{
				Name:        check.Name,
				Description: check.Description,
				Detail:      line,
				Source:      check.Name,
				Location:    file,
			}
			log.Errorf("Finding: %v", finding)
			findings = append(findings, finding)
		}
	}

	//log.Debugf("Check: %s", check.Magic)
	matched := r.MatchString(line)
	if matched {
		finding := finding.Finding{
			Name:        check.Name,
			Description: check.Description,
			Detail:      line,
			Source:      check.Name,
			Location:    file,
		}
		log.Errorf("Finding: %v", finding)
		findings = append(findings, finding)
	}

	return findings
}

func getAllChecks(opts options.Options) []check.Check {
	var checks []check.Check
	checks = append(checks, getChecks("check/injections.json")...)
	checks = append(checks, getChecks("check/secrets.json")...)
	checks = append(checks, getChecks("check/files.json")...)
	checks = append(checks, getChecks("check/unescaped.json")...)
	return checks
}

func getChecks(file string) []check.Check {
	var checks []check.Check
	rfile, err := os.Open(file)
	if err != nil {
		log.Error(err)
	}
	bytes, err := ioutil.ReadAll(rfile)
	if err != nil {
		log.Error(err)
	}
	json.Unmarshal(bytes, &checks)
	return checks
}

func getFiles(options options.Options) []string {
	var files []string
	err := filepath.Walk(options.Directory,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.IsDir() {
				//log.Debugf("File %s", path)
				files = append(files, path)
			}
			return nil
		})
	if err != nil {
		log.Error(err)
	}
	return files
}

func readLines(path string) ([]string, error) {
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

func buildExamineOptions(cmd *cobra.Command) options.Options {
	directory := viper.GetString("directory")
	options := options.Options{
		Directory: directory,
	}

	debug := viper.GetBool("debug")
	if debug != true {
		log.SetLevel(log.InfoLevel)
	}
	log.Debug("Captured options: ")
	log.Debug(options)

	return options
}

func init() {
	rootCmd.AddCommand(examineCmd)

	examineCmd.PersistentFlags().String("directory", "", "The directory to examine.")
	examineCmd.MarkFlagRequired("directory")

	viper.BindPFlag("directory", examineCmd.PersistentFlags().Lookup("directory"))

	log.SetFormatter(&log.TextFormatter{})
	log.SetLevel(log.DebugLevel)

}
