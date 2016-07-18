// Copyright © 2016 NAME HERE <EMAIL ADDRESS>
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
	"fmt"
	"net/http"
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/lytics/cloudstorage"
	"github.com/spf13/cobra"
	"google.golang.org/cloud/compute/metadata"
)

var gceproject string
var gc *http.Client
var topic string
var keyPath string

func GCS(projectid string) cloudstorage.GoogleOAuthClient {
	onGce := metadata.OnGCE()
	gcsctx := &cloudstorage.CloudStoreContext{
		LogggingContext: "secure-config",
		TokenSource:     cloudstorage.GCEDefaultOAuthToken,
		Project:         projectid,
		Bucket:          "neh",
	}
	if onGce {
		gcsctx.TokenSource = cloudstorage.GCEMetaKeySource
	}

	// Create http client with Google context auth
	googleClient, err := cloudstorage.NewGoogleClient(gcsctx)
	if err != nil {
		log.Errorf("failed to create google storage Client: %v", err)
		os.Exit(1)
	}

	return googleClient
}

func initClient() *http.Client {
	metaproject, _ := metadata.ProjectID()
	if gceproject == "" && metaproject == "" {
		log.Errorf("No project specified")
		os.Exit(1)
	} else if gceproject == "" && metaproject != "" {
		gceproject = metaproject
	}

	gcs := GCS(gceproject)
	gc := gcs.Client()
	log.Debugf("Google Auth: %#v", gc)
	return gc
}

func init() {
	log.SetLevel(log.DebugLevel)

	// Here you will define your flags and configuration settings.
	// Cobra supports Persistent Flags, which, if defined here,
	// will be global for your application.
	RootCmd.PersistentFlags().StringVar(&gceproject, "project", "", "GCE Project")
	RootCmd.PersistentFlags().StringVar(&topic, "topic", "", "PubSub topic")
	RootCmd.PersistentFlags().StringVar(&keyPath, "key", "", "PubSub service account key path")

	//cobra.OnInitialize(initClient)
}

// This represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "pubbing",
	Short: "Google PubSub test framework",
	Long:  ``,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) { log.Infof("pubbing called without command") },
}

// Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

/*
// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" { // enable ability to specify config file via flag
		viper.SetConfigFile(cfgFile)
	}

	viper.SetConfigName(".pubbing") // name of config file (without extension)
	viper.AddConfigPath("$HOME")    // adding home directory as first search path
	viper.AutomaticEnv()            // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
*/
