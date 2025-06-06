# Logging with Loki, Promtail, and Grafana

This document describes how to set up and use the logging infrastructure for the Book System application.

## Architecture

- **slog**: Structured logging in Go
- **Loki**: Log aggregation system
- **Promtail**: Log collector that ships logs to Loki
- **Grafana**: Visualization platform for logs

## Setup

1. Start the services:
   ```bash
   docker-compose up -d loki promtail grafana
   ```

2. Access the services:
   - Grafana: http://localhost:3000 (admin/admin)
   - Loki: http://localhost:3100
   - Promtail: http://localhost:9080

## Configuring Grafana

1. Log in to Grafana (default credentials: admin/admin)
2. Add Loki as a data source:
   - Go to Configuration > Data Sources
   - Click "Add data source"
   - Select "Loki"
   - Set URL to `http://loki:3100`
   - Click "Save & Test"

## Querying Logs in Grafana

1. Go to Explore in the left sidebar
2. Select the Loki data source
3. Use LogQL to query logs, for example:
   ```
   {container="book_system"} | json
   ```

## Log Levels

- `DEBUG`: Detailed information for debugging
- `INFO`: General operational information
- `WARN`: Non-critical issues that should be addressed
- `ERROR`: Critical errors that need attention

## Adding Context to Logs

Use the `With` method to add context to your logs:

```go
log := logger.New(slog.LevelInfo)
log = log.With(
    "request_id", requestID,
    "user_id", userID,
)
log.Info("User logged in")
```

## Best Practices

1. Use structured logging with key-value pairs
2. Include relevant context in log messages
3. Use appropriate log levels
4. Avoid logging sensitive information
5. Use the `With` method to add common fields to all log messages
