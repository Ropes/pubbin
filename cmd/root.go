// Copyright © 2016 Josh Roppo joshroppo@gmail.com
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

var (
	GC         *http.Client
	Gceproject string
	Topic      string
	KeyPath    string
	Loglvl     string
	Logfmt     string
)

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

func logsetup() {
	lvl, err := log.ParseLevel(Loglvl)
	if err != nil {
		log.Errorf("error parsing loglevel: %v", err)
	}
	log.SetLevel(lvl)

	log.Debugf("Logfmt: %s", Logfmt)
	if Logfmt == "json" {
		log.SetFormatter(&log.JSONFormatter{})
	}
}

func initClient() *http.Client {
	metaproject, _ := metadata.ProjectID()
	if Gceproject == "" && metaproject == "" {
		log.Errorf("No project specified")
		os.Exit(1)
	} else if Gceproject == "" && metaproject != "" {
		Gceproject = metaproject
	}

	gcs := GCS(Gceproject)
	gc := gcs.Client()
	log.Debugf("Google Auth: %#v", gc)
	return gc
}

func rootinit(Cmd *cobra.Command) {
	RootCmd.PersistentFlags().StringVar(&Gceproject, "project", "", "GCE Project")
	RootCmd.PersistentFlags().StringVar(&Topic, "topic", "", "PubSub topic")
	RootCmd.PersistentFlags().StringVar(&KeyPath, "key", "", "PubSub service account key path")
	RootCmd.PersistentFlags().StringVar(&Loglvl, "log", "info", "logging level; debug,info,warn,error")
	RootCmd.PersistentFlags().StringVar(&Logfmt, "logfmt", "text", "logging format: text,json")
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

func init() {
	rootinit(RootCmd)
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// pubCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// pubCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}

// Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}
