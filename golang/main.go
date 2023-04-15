package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/K-Phoen/grabana"
	"github.com/K-Phoen/grabana/dashboard"
	"github.com/K-Phoen/grabana/row"
	"github.com/K-Phoen/grabana/stat"
	"github.com/K-Phoen/grabana/target/prometheus"
	"github.com/K-Phoen/grabana/timeseries"
	"github.com/K-Phoen/grabana/timeseries/axis"
	"github.com/K-Phoen/grabana/variable/interval"
)

func main() {
	grafanaServer := os.Getenv("GRAFANA_SERVER")
	grafanaApiKey := os.Getenv("GRAFANA_API_KEY")

	ctx := context.Background()
	client := grabana.NewClient(&http.Client{}, grafanaServer, grabana.WithAPIToken(grafanaApiKey))

	// create the folder holding the dashboard for the service
	folder, err := client.FindOrCreateFolder(ctx, "Grabana - Golang")
	if err != nil {
		fmt.Printf("Could not find or create folder: %s\n", err)
		os.Exit(1)
	}

	builder, err := dashboard.New(
		"Kubernetes Cluster",
		dashboard.UID("k8s-cluster-golang"),
		dashboard.AutoRefresh("30s"),
		dashboard.Time("now-30m", "now"),
		dashboard.Tags([]string{"generated"}),
		dashboard.VariableAsInterval(
			"interval",
			interval.Values([]string{"30s", "1m", "5m", "10m", "30m", "1h", "6h", "12h"}),
			interval.Default("5m"),
		),
		dashboard.Row(
			"Overview",

			row.WithTimeSeries(
				"CPU Utilization by namespace",
				timeseries.DataSource("Prometheus"),
				timeseries.WithPrometheusTarget("sum(rate(container_cpu_usage_seconds_total{image!=\"\"}[$__rate_interval])) by (namespace)",
					prometheus.Legend("{{ namespace }}")),
				timeseries.Height("400px"),
			),

			row.WithTimeSeries(
				"Memory Utilization by namespace",
				timeseries.DataSource("Prometheus"),
				timeseries.WithPrometheusTarget("sum(container_memory_working_set_bytes{image!=\"\"}) by (namespace)",
					prometheus.Legend("{{ namespace }}")),
				timeseries.Axis(
					axis.Unit("bytes"),
					axis.Label("Memory"),
					axis.SoftMin(0),
				),
				timeseries.Legend(timeseries.Last, timeseries.AsTable),
				timeseries.Height("400px"),
			),

			row.WithStat("Running Pods",
				stat.WithPrometheusTarget("sum(kube_pod_status_phase{phase='Running'})"),
				stat.Height("400px"),
			),
		),
	)
	if err != nil {
		fmt.Printf("Could not build dashboard: %s\n", err)
		os.Exit(1)
	}

	_, err = client.UpsertDashboard(ctx, folder, builder)
	if err != nil {
		fmt.Printf("Could not create dashboard: %s\n", err)
		os.Exit(1)
	}
}
