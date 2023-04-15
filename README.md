# GitOps for Grafana dashboards

Create Grafana dashboard in GitOps way.

## Install dashboards

1. Export grafana API key and server URL.
```
export GRAFANA_SERVER=
export GRAFANA_API_KEY=
```

2. Create dashboard from Golang code:
```
go run golang/main.go
```

3. Create dashboard from YAML:
```
go run yaml/main.go
```
