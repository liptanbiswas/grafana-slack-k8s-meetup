package main

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/K-Phoen/grabana"
	"github.com/K-Phoen/grabana/decoder"
)

func main() {

	content, err := os.ReadFile("yaml/dashboard.yaml")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not read file: %s\n", err)
		os.Exit(1)
	}

	dashboard, err := decoder.UnmarshalYAML(bytes.NewBuffer(content))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not parse file: %s\n", err)
		os.Exit(1)
	}

	grafanaServer := os.Getenv("GRAFANA_SERVER")
	grafanaApiKey := os.Getenv("GRAFANA_API_KEY")

	ctx := context.Background()
	client := grabana.NewClient(&http.Client{}, grafanaServer, grabana.WithAPIToken(grafanaApiKey))

	// create the folder holding the dashboard for the service
	folder, err := client.FindOrCreateFolder(ctx, "Grabana - Yaml")
	if err != nil {
		fmt.Printf("Could not find or create folder: %s\n", err)
		os.Exit(1)
	}

	if _, err := client.UpsertDashboard(ctx, folder, dashboard); err != nil {
		fmt.Printf("Could not create dashboard: %s\n", err)
		os.Exit(1)
	}
}
