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
	"time"

	"golang.org/x/net/context"

	log "github.com/Sirupsen/logrus"
	"github.com/spf13/cobra"
	"google.golang.org/cloud"
	"google.golang.org/cloud/pubsub"
)

// pubCmd represents the pub command
var pubCmd = &cobra.Command{
	Use:   "pub",
	Short: "publish messages to defined topic",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Infof("pub called on topic: %s", Topic)

		if Gceproject == "" || Topic == "" {
			log.Errorf("GCE project and topic must be defined")
			os.Exit(1)
		}
		ctx := context.Background()
		pubsubClient := initClient()
		gctx := cloud.NewContext(Gceproject, pubsubClient)

		var psClient *pubsub.Client
		if KeyPath != "" {
			psClient = JWTClientInit(&ctx)
		} else {
			psClient = GCEClientInit(&ctx, Gceproject)
		}
		if psClient == nil {
			log.Errorf("PubSub client is nil")
			os.Exit(1)
		}

		topic := psClient.Topic(Topic)
		bytes := []byte(fmt.Sprintf("helloworld %v", time.Now()))
		ids, err := topic.Publish(gctx, &pubsub.Message{Data: bytes})
		if err != nil {
			log.Errorf("error publishing messages: %v", err)
			os.Exit(1)
		}
		for _, id := range ids {
			log.Infof("%#v", id)
		}
	},
}

func init() {
	RootCmd.AddCommand(pubCmd)
}
