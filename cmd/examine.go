// Copyright © 2019-2020 Matt Konda <mkonda@jemurai.com>
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
	"strconv"
	"sync"
	"time"

	"github.com/jemurai/crush/check"
	"github.com/jemurai/crush/options"
	"github.com/jemurai/crush/utils"

	"github.com/jemurai/fkit/finding"
	fkitutils "github.com/jemurai/fkit/utils"

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
		log.Debugf("Examine command on: %s", opts.Directory)
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

		if opts.Compare != "" {
			added := fkitutils.CompareFileAndArray(opts.Compare, findings)
			printFindings(added, checks, files) // only printing new added findings
		} else {
			printFindings(findings, checks, files)
		}
		utils.Timing(start, "Elasped time: %f")
	},
}

func printFindings(findings []finding.Finding, checks []check.Check, files []string) {
	if len(findings) > 0 {
		fjson, _ := json.MarshalIndent(findings, "", " ")
		fmt.Printf("%s", fjson)
		howMany := check.CountHowManyRan(checks)
		log.Debugf("\n\nSummary:\n\tFiles processed: %v \n\tChecks: %v (Not working yet) \n\tIssues: %v", len(files), howMany, len(findings))
	} else {
		fmt.Print("[]")
	}
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
			if check.AppliesToTag(checks[j], options.Tag) &&
				check.AppliesToExt(checks[j], filepath.Ext(fn), options.Ext) &&
				check.AppliesBasedOnThreshold(checks[j], options.Threshold) &&
				check.AppliesBasedOnComment(lines[i], filepath.Ext(fn)) {
				fs := doCheck(fn, i, checks[j], lines[i], options)
				findings = append(findings, fs...)
			}
		}
	}
	// log.Debugf("\tProcessed %s", fn)
	return findings
}

func doCheck(file string, lineno int, check check.Check, line string, options options.Options) []finding.Finding {
	var findings []finding.Finding

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
	tag := viper.GetString("tag")
	ext := viper.GetString("ext")
	compare := viper.GetString("compare")
	threshold := viper.GetFloat64("threshold")

	options := options.Options{
		Directory: directory,
		Tag:       tag,
		Ext:       ext,
		Compare:   compare,
		Threshold: threshold,
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

	examineCmd.PersistentFlags().String("tag", "", "The tag for checks to run.")
	viper.BindPFlag("tag", examineCmd.PersistentFlags().Lookup("tag"))

	examineCmd.PersistentFlags().String("ext", "", "The file extension for checks to run.")
	viper.BindPFlag("ext", examineCmd.PersistentFlags().Lookup("ext"))

	examineCmd.PersistentFlags().String("compare", "", "The file to compare new results to.")
	viper.BindPFlag("compare", examineCmd.PersistentFlags().Lookup("compare"))

	examineCmd.PersistentFlags().Float64("threshold", 5.0, "The threshold of confidence we want to hold findings to.")
	viper.BindPFlag("threshold", examineCmd.PersistentFlags().Lookup("threshold"))

	log.SetFormatter(&log.TextFormatter{})
	log.SetLevel(log.DebugLevel)

}
