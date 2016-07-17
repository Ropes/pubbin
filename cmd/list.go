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
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/spf13/cobra"
	"golang.org/x/net/context"
	"google.golang.org/api/pubsub/v1"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Info("list called")
		gc := initClient()

		log.Debugf("google client: %#v", gc)
		c, err := pubsub.New(gc)
		if err != nil {
			log.Errorf("pubsub client connection error: %v", err)
			os.Exit(1)
		}

		// The name of the cloud project that topics belong to.
		project := fmt.Sprintf("projects/%s", gceproject)
		ctx := context.Background()

		call := c.Projects.Topics.List(project)
		if err := call.Pages(ctx, func(page *pubsub.ListTopicsResponse) error {
			for _, v := range page.Topics {
				// TODO: Use v.
				//_ = v
				log.Info("Topic found: %s", v)
			}
			return nil // NOTE: returning a non-nil error stops pagination.
		}); err != nil {
			log.Errorf("No Topics returned")
			// TODO: Handle error.
		}

	},
}

func init() {
	RootCmd.AddCommand(listCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// listCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// listCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}
