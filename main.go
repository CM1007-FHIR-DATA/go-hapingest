package main

import (
	"context"
	"fmt"
	"log"
	"path/filepath"
	"sync"
	"time"

	"github.com/CM1007-FHIR-DATA/go-hapingest/config"
	"github.com/CM1007-FHIR-DATA/go-hapingest/internal/util"
	"github.com/CM1007-FHIR-DATA/go-hapingest/pkg/fhir"
	"github.com/CM1007-FHIR-DATA/go-hapingest/pkg/files/server"
	"github.com/CM1007-FHIR-DATA/go-hapingest/pkg/manifest"
)

func main() {
	cfg := config.GetInstance()
	fmt.Println("Current Configuration: ")
	fmt.Println(cfg.String())

	var wg sync.WaitGroup
	ctx, cancelFunc := context.WithCancel(context.Background())
	defer func() {
		cancelFunc()
		wg.Wait()
	}()
	server := server.NewServer(cfg.DataDir, cfg.Port)
	go server.Start(ctx, &wg)

	files, err := filepath.Glob(filepath.Join(cfg.DataDir, "*.ndjson"))
	if err != nil {
		log.Fatalf("Error reading directory %s: %v", cfg.DataDir, err)
	}
	if len(files) == 0 {
		log.Fatalf("No .ndjson files found in %s", cfg.DataDir)
	}

	if cfg.PingServer {
		fhir.PingServer(cfg.FHIRServerURL)
	}

	manifestData := manifest.CreateManifest(files, cfg.URLBase+":"+cfg.Port)
	for _, manifest := range manifestData {
		formattedData, err := util.FormatJSON(manifest)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(string(formattedData) + "\n\n")
		contentLocation := fhir.PostManifest(cfg.FHIRServerURL, manifest)
		if contentLocation != "" {
			success := fhir.CheckJobStatus(contentLocation)
			if !success {
				log.Fatalf("Error checking job status: %v\n", err)
			} else {
				fmt.Println("Job completed successfully.")
			}
		}
		time.Sleep(1 * time.Second)
	}
}
