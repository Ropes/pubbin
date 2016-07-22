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
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"golang.org/x/net/context"
	"google.golang.org/cloud"
	"google.golang.org/cloud/pubsub"

	log "github.com/Sirupsen/logrus"
	"github.com/schmichael/legendarygopher/lg"
	"github.com/spf13/cobra"
)

var (
	num          int
	batchInt     int
	httpendpoint string
)

// gopherpumpCmd represents the gopherpump command
var gopherpumpCmd = &cobra.Command{
	Use:   "gopherpump",
	Short: "Read LegendaryGopher api and pump Figures into PubSub",
	Long:  `Queries a running instances of LegendaryGopher API; marshals the messages into PubSub`,
	Run: func(cmd *cobra.Command, args []string) {
		logsetup()
		log.Debugf("gopherpump called: num: %d batch: %d http: %s", num, batchInt, httpendpoint)
		// Listen for kill signals
		quit := make(chan os.Signal, 1)
		figureChan := make(chan *lg.Figure)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

		// ------------------------------------------------------------------
		// Query LegendaryGopher API
		resp, err := http.Get(httpendpoint)
		if err != nil {
			log.Errorf("error reading http endpoint: '%s'  %v", httpendpoint, err)
			os.Exit(1)
		}
		figureBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Errorf("error reading body: %v", err)
			os.Exit(1)
		}
		log.Debugf("figure bytes read %d", len(figureBytes))

		// Unmarshal data into list of Figures
		figures := make([]*lg.Figure, 0)
		defer close(figureChan)
		err = json.Unmarshal(figureBytes, &figures)
		if err != nil {
			log.Errorf("error unmarshaling figures: %v", err)
			os.Exit(1)
		}
		log.Debugf("Unmarshaling complete")

		// Pump Figures from
		go func() {
			for i, f := range figures {
				log.Debugf("[%d]: %#v", i, *f)
				figureChan <- f
			}
		}()

		// ------------------------------------------------------------------
		// PubSub Pump
		ctx := context.Background()
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

		if Gceproject == "" || Topic == "" {
			log.Errorf("GCE project and topic must be defined")
			os.Exit(1)
		}
		psHTTPClient := initClient()
		gctx := cloud.NewContext(Gceproject, psHTTPClient)
		topic := psClient.Topic(Topic)

		// Write Figures to PubSub
		i := 0
		exit := false
		msgs := make([]*pubsub.Message, 0)
		for exit == false && i < num {
			select {
			case f := <-figureChan:
				i++
				log.Debugf("figure: %#v", *f)

				attrs := map[string]string{"race": f.Race, "caste": f.Caste, "name": f.Name}

				if len(msgs) < batchInt {
					m := &pubsub.Message{Attributes: attrs, Data: []byte(f.Name)}
					msgs = append(msgs, m)
				} else {
					log.Infof("Publishing %d to %s", len(msgs), Topic)
					ids, err := topic.Publish(gctx, msgs...)
					if err != nil {
						log.Errorf("error publishing: %v", err)
					}
					log.Debugf("Message IDs\n%#v", ids)
					msgs = nil
				}

			case <-time.After(time.Second * 1):
				log.Debugf("publisher heartbeat")
			case <-quit:
				log.Warnf("quit signal caught")
				exit = true
			}
		}
		log.Infof("Figures: %d", i)
	},
}

func init() {
	RootCmd.AddCommand(gopherpumpCmd)

	gopherpumpCmd.Flags().IntVar(&num, "num", 100, "Number of entities to write to pubsub")
	gopherpumpCmd.Flags().IntVar(&batchInt, "batch", 50, "PubSub publishing batch sizes")
	gopherpumpCmd.Flags().StringVar(&httpendpoint, "figures", "http://localhost:6565/api/figures", "JSON API endpoint to read figure definitions from")
}
