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
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/spf13/cobra"
	"golang.org/x/net/context"
	"google.golang.org/api/pubsub/v1"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all available topics in a project",
	Long:  `List all available topics in a project`,
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
				log.Infof("Topic found: %s", v.Name)
			}
			return nil // NOTE: returning a non-nil error stops pagination.
		}); err != nil {
			log.Errorf("No Topics returned")
		}
	},
}

func init() {
	RootCmd.AddCommand(listCmd)
}
