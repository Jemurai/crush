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

package cmd

import (
	"os"

	homedir "github.com/mitchellh/go-homedir"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string
var debug bool
var tag string
var ext string
var compare string
var threshold float64

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "crush",
	Short: "Code review helps us.",
	Long: `Find potential security issues by looking at code.
	
Many types of security issues cannot be found with static 
analysis and are amenable to developer review.

Use crush help to get more information about using crush.
`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Error(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.crush.yaml)")
	rootCmd.PersistentFlags().BoolVar(&debug, "debug", false, "debug mode")
	viper.BindPFlag("debug", rootCmd.PersistentFlags().Lookup("debug"))

	rootCmd.PersistentFlags().String("tag", "", "The tag for checks to run.")
	viper.BindPFlag("tag", rootCmd.PersistentFlags().Lookup("tag"))

	rootCmd.PersistentFlags().StringVar(&ext, "ext", "", "The file extension for checks to run.")
	viper.BindPFlag("ext", rootCmd.PersistentFlags().Lookup("ext"))

	rootCmd.PersistentFlags().String("compare", "", "The file to compare new results to.")
	viper.BindPFlag("compare", rootCmd.PersistentFlags().Lookup("compare"))

	rootCmd.PersistentFlags().Float64("threshold", 5.0, "The threshold of confidence we want to hold findings to.")
	viper.BindPFlag("threshold", rootCmd.PersistentFlags().Lookup("threshold"))

	log.SetFormatter(&log.TextFormatter{})
	log.SetLevel(log.DebugLevel)
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if debug {
		log.SetLevel(log.DebugLevel)
	} else {
		log.SetLevel(log.InfoLevel)
	}

	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			log.Error(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".crh" (with extension!!!).
		viper.AddConfigPath(home)
		viper.SetConfigName(".crush")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		log.Debug("Using config file:", viper.ConfigFileUsed())
	} else {
		//Uncomment if problems picking up config file.
		//fmt.Println(err)
	}
}
