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
	"io/ioutil"
	"os"
	"os/signal"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/spf13/cobra"
	"golang.org/x/net/context"
	"golang.org/x/oauth2/google"
	"google.golang.org/cloud"
	"google.golang.org/cloud/pubsub"
)

var subscription string
var numConsume int

var quit chan os.Signal

func shouldQuit() bool {
	select {
	case <-quit:
		signal.Stop(quit)
		close(quit)
		return true
	default:
		return false
	}
}

// subCmd represents the sub command
var subCmd = &cobra.Command{
	Use:   "sub",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Infof("sub called on topic: %s", topic)

		quit := make(chan os.Signal, 1)
		signal.Notify(quit, os.Interrupt)

		if gceproject == "" || topic == "" || subscription == "" || keyPath == "" {
			log.Errorf("GCE project, subscription, keypath, and topic must be defined")
			os.Exit(1)
		}

		jsonKey, err := ioutil.ReadFile(keyPath)
		if err != nil {
			log.Errorf("error reading keyfile: %v", err)
			os.Exit(1)
		}

		conf, err := google.JWTConfigFromJSON(jsonKey, pubsub.ScopePubSub)
		if err != nil {
			log.Errorf("error creating conf file: %v", err)
		}

		ctx := context.Background()
		oauthTokenSource := conf.TokenSource(ctx)

		psClient, err := pubsub.NewClient(ctx, gceproject, cloud.WithTokenSource(oauthTokenSource))
		if err != nil {
			log.Errorf("error creating pubsub client: %v", err)
			os.Exit(1)
		}
		log.Debugf("client: %#v", psClient)

		sub := psClient.Subscription(subscription)

		it, err := sub.Pull(ctx, pubsub.MaxExtension(time.Minute))
		if err != nil {
			log.Errorf("error creating pubsub iterator: %v", err)
		}
		defer it.Stop()

		for !shouldQuit() {
			m, err := it.Next()
			if err != nil {
				log.Errorf("error reading from iterator: %v", err)
			}
			log.Infof("msg read: %#v", m)
		}

		/*
			gc := initClient()
			ctx := context.Background()
			gctx := cloud.NewContext(gceproject, gc)
			psClient, err := pubsub.NewClient(ctx, gceproject, gctx)
			if err != nil {
				log.Errorf("error creating pubsub.Client: %v", err)
				os.Exit(1)
			}

			sub := psClient.Subscription(subscription)

			it, err := sub.Pull(gctx, pubsub.MaxExtension(time.Minute))
			if err != nil {
				log.Errorf("error constructing iterator: %v", err)
				os.Exit(1)
			}

			defer it.Stop()
			go func() {
				<-quit
				it.Stop()
			}()

			for i := 0; i < *numConsume; i++ {
				m, err := it.Next()
				if err == pubsub.Done {
					break
				}
				if err != nil {
					log.Errorf("advancing iterator %v", err)
					break
				}
				log.Infof("message: %#v", m)
				m.Done(true)
			}

				msgs, err := pubsub.Pull(gctx, topic, 10)
				if err != nil {
					log.Errorf("error pulling messages: %v", err)
				}

				for _, m := range msgs {
					//TODO: ACK messages
					log.Infof("   msg: %#v", m)

				}
		*/
	},
}

func init() {
	RootCmd.AddCommand(subCmd)
	RootCmd.PersistentFlags().StringVar(&subscription, "sub", "", "PubSub subscription")
	RootCmd.PersistentFlags().IntVar(&numConsume, "num", 10, "Messages to consume")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// subCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// subCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}
