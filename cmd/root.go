// Copyright © 2017 Farhad Farahi
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
	"fmt"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
	"os"
	"strings"
	"utmpsave/utmp"
)

var (
	cfgFile     string
	source      string
	destination string
)

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "utmpsave",
	Short: "Dumps linux lastlog into a human readable file",
	Long: `Dumps linux lastlog into a human readable file.
Can be scheduled, utmpsave will not store duplicate records.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: rootCmd,
}

func rootCmd(cmd *cobra.Command, args []string) {
	if versionFlag := getFlagBoolPtr(cmd, "version"); versionFlag != nil {
		fmt.Println("utmpsave v1.0.0")
	} else {
		utmp.UtmpSave(source, destination)
	}
}

func getFlagBoolPtr(cmd *cobra.Command, flag string) *bool {
	f := cmd.Flags().Lookup(flag)
	if f == nil {
		log.Printf("Flag accessed but not defined for command %s: %s", cmd.Name(), flag)
	}
	// Check if flag was not set at all.
	if !f.Changed && f.DefValue == f.Value.String() {
		return nil
	}
	var ret bool
	// Caseless compare.
	if strings.ToLower(f.Value.String()) == "true" {
		ret = true
	} else {
		ret = false
	}
	return &ret
}

// Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	//cobra.OnInitialize(initConfig)
	cobra.OnInitialize()
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	//RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.utmpdump.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	RootCmd.Flags().BoolP("version", "v", false, "Prints version info")
	RootCmd.Flags().StringVarP(&source, "source", "s", "/var/log/wtmp", "Source wtmp file path")
	RootCmd.Flags().StringVarP(&destination, "destination", "d", "/tmp/utmpsave.dmp", "Destination dump file path")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".utmpdump" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".utmpsave")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
