package main

import (
	"flag"
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/lytics/cloudstorage"
	"google.golang.org/cloud/compute/metadata"
)

var gceproject string

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

func main() {
	log.SetLevel(log.DebugLevel)

	flag.StringVar(&gceproject, "project", "", "GCE project to use")
	metaproject, _ := metadata.ProjectID()

	flag.Parse()

	if gceproject == "" && metaproject == "" {
		log.Errorf("No project specified")
		os.Exit(1)
	} else if gceproject == "" && metaproject != "" {
		gceproject = metaproject
	}

	gcs := GCS(gceproject)
	gc := *gcs.Client()
	log.Debugf("Google Auth: %#v", gc)

}
