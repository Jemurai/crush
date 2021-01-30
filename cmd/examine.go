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
	Short: "Examine a directory - and optionally run all the checks",
	Long: `Examine all the source code in a directory.
	
Behind the scenes, crush checks a variety of things.
See the checks in /check and .json files.`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		start := time.Now()
		opts := buildExamineOptions(cmd, args)
		log.Debugf("Examine command on: %s", opts.Directory)
		checks := getAllChecks(opts)
		files := check.GetFiles(opts)
		var findings []finding.Finding
		var wg sync.WaitGroup
		for i := 0; i < len(files); i++ {
			wg.Add(1)
			fn := files[i]
			go func(fn string, opts options.Options) {
				defer wg.Done()
				findings = append(findings, check.ProcessFile(fn, checks, opts)...)
			}(fn, opts)
		}
		wg.Wait()

		if opts.Compare != "" {
			added := fkitutils.CompareFileAndArray(opts.Compare, findings)
			check.PrintFindings(added, checks, files) // only printing new added findings
		} else {
			check.PrintFindings(findings, checks, files)
		}
		utils.Timing(start, "Elasped time: %f")
		log.Debug(opts)

	},
}

func getAllChecks(opts options.Options) []check.Check {
	var checks []check.Check
	checks = append(checks, check.GetChecks("checks/injections.json")...)
	checks = append(checks, check.GetChecks("checks/secrets.json")...)
	checks = append(checks, check.GetChecks("checks/files.json")...)
	checks = append(checks, check.GetChecks("checks/unescaped.json")...)
	checks = append(checks, check.GetChecks("checks/mobile.json")...)
	checks = append(checks, check.GetChecks("checks/badwords.json")...)
	return checks
}

func buildExamineOptions(cmd *cobra.Command, args []string) options.Options {
	directory := args[0]
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
}
