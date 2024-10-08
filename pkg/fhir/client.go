package fhir

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

func PostManifest(fhirServerURL string, manifest interface{}) string {
	jsonData, err := json.MarshalIndent(manifest, "", "  ")
	if err != nil {
		log.Fatalf("Error marshalling manifest: %v", err)
	}

	url := fmt.Sprintf("%s/$import", fhirServerURL)
	req, err := http.NewRequest("POST", url, bytes.NewReader(jsonData))
	if err != nil {
		log.Fatalf("Error creating request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/fhir+json")
	req.Header.Set("Prefer", "respond-async")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Error posting manifest: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusAccepted {
		fmt.Println("Import request accepted.")
		contentLocation := resp.Header.Get("Content-Location")
		if contentLocation != "" {
			fmt.Printf("Check status at: %s\n", contentLocation)
			return contentLocation
		} else {
			fmt.Println("No Content-Location header found.")
		}
	} else {
		fmt.Printf("Failed to post manifest. Status code: %d\n", resp.StatusCode)
		bodyBytes, _ := io.ReadAll(resp.Body)
		fmt.Println(string(bodyBytes))
	}
	return ""
}

func CheckJobStatus(contentLocation string) bool {
	if contentLocation == "" {
		fmt.Println("No valid content location provided.")
		return false
	}

	client := &http.Client{}
	req, err := http.NewRequest("GET", contentLocation, nil)
	if err != nil {
		log.Fatalf("Error creating request to check job status: %v", err)
	}

	req.Header.Set("Accept", "application/fhir+json")

	for {
		resp, err := client.Do(req)
		if err != nil {
			log.Fatalf("Error checking job status: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode == http.StatusOK {
			fmt.Println("Job completed successfully.")
			bodyBytes, _ := io.ReadAll(resp.Body)
			fmt.Println("Response:\n", string(bodyBytes))
			return true
		} else if resp.StatusCode == http.StatusAccepted {
			fmt.Println("Job is still in progress...")
		} else {
			fmt.Printf("Failed to get job status. Status code: %d\n", resp.StatusCode)
			bodyBytes, _ := io.ReadAll(resp.Body)
			fmt.Println(string(bodyBytes))
			return false
		}

		time.Sleep(5 * time.Second)
	}
}

// PingServer continuously pings the HAPI FHIR server until it is up and responding.
// it is blocking and makes sure the HAPI FHIR server is up and running before it starts.
// HAPI can take a long time to start sometimes so this is here to handle that
func PingServer(fhirServerURL string) {
	client := &http.Client{}
	url := fmt.Sprintf("%s/$meta", fhirServerURL)

	for {
		resp, err := client.Get(url)
		if err != nil {
			fmt.Printf("Error pinging server: %v. Retrying...\n", err)
		} else {
			defer resp.Body.Close()
			if resp.StatusCode == http.StatusOK {
				fmt.Println("FHIR server is up and running.")
				return
			} else {
				fmt.Printf("Server responded with status code: %d. Retrying...\n", resp.StatusCode)
			}
		}

		time.Sleep(5 * time.Second)
	}
}
