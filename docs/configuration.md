# GoProxy Configuration Guide

GoProxy uses a YAML configuration file to customize its behavior. This document provides a detailed explanation of all available configuration options.

## Configuration File Structure

The configuration file is divided into several sections, each controlling a different aspect of GoProxy's behavior:

- Server
- Proxy
- Load Balancing
- TLS
- Logging
- Metrics
- Rate Limiting
- Caching

## Server Settings

These settings control the basic behavior of the GoProxy server.

```yaml
server:
  listen_addr: ":8080"
  read_timeout: 5
  write_timeout: 10
  idle_timeout: 120
```

- `listen_addr`: The address and port on which GoProxy will listen for incoming requests. Format is `"host:port"`. Use `:port` to listen on all interfaces.
- `read_timeout`: Maximum duration (in seconds) for reading the entire request, including the body.
- `write_timeout`: Maximum duration (in seconds) before timing out writes of the response.
- `idle_timeout`: Maximum amount of time (in seconds) to wait for the next request when keep-alives are enabled.

## Proxy Settings

These settings configure how GoProxy forwards requests to the target server.

```yaml
proxy:
  target_addr: "http://localhost:8000"
  max_idle_conns: 100
  dial_timeout: 10
```

- `target_addr`: The address of the backend server to which GoProxy will forward requests.
- `max_idle_conns`: The maximum number of idle (keep-alive) connections between the proxy and the backend.
- `dial_timeout`: The maximum amount of time (in seconds) to wait for a connection to the backend.

## Load Balancing Settings

(Note: This feature is planned for future implementation)

```yaml
load_balancing:
  enabled: false
  algorithm: "round_robin"
  backends: []
```

- `enabled`: Set to `true` to enable load balancing.
- `algorithm`: The load balancing algorithm to use. Options will include "round_robin", "least_connections", etc.
- `backends`: A list of backend server addresses for load balancing.

## TLS Settings

(Note: This feature is planned for future implementation)

```yaml
tls:
  enabled: false
  cert_file: ""
  key_file: ""
```

- `enabled`: Set to `true` to enable TLS.
- `cert_file`: Path to the TLS certificate file.
- `key_file`: Path to the TLS private key file.

## Logging Settings

Configure the logging behavior of GoProxy.

```yaml
logging:
  level: "info"
  format: "json"
```

- `level`: The minimum log level to output. Options are "debug", "info", "warn", "error".
- `format`: The format of log output. Options are "json" or "text".

## Metrics Settings

(Note: This feature is planned for future implementation)

```yaml
metrics:
  enabled: false
  address: ":9090"
```

- `enabled`: Set to `true` to enable Prometheus metrics.
- `address`: The address on which to expose the Prometheus metrics.

## Rate Limiting Settings

(Note: This feature is planned for future implementation)

```yaml
rate_limiting:
  enabled: false
  requests_per_second: 100
  burst: 50
```

- `enabled`: Set to `true` to enable rate limiting.
- `requests_per_second`: The number of requests allowed per second.
- `burst`: The maximum number of requests allowed to exceed the rate in a short burst.

## Caching Settings

(Note: This feature is planned for future implementation)

```yaml
caching:
  enabled: false
  default_ttl: 300
  max_size_mb: 100
```

- `enabled`: Set to `true` to enable caching.
- `default_ttl`: The default time-to-live for cached items, in seconds.
- `max_size_mb`: The maximum size of the cache, in megabytes.

## Environment Variable Overrides

GoProxy allows overriding configuration settings using environment variables. The format for environment variables is:

```
GOPROXY_SECTION_KEY=value
```

For example, to override the `listen_addr` in the `server` section:

```
GOPROXY_SERVER_LISTEN_ADDR=:9090
```

Environment variables take precedence over values in the configuration file.

## Sample Complete Configuration

Here's a sample configuration file with all options:

```yaml
server:
  listen_addr: ":8080"
  read_timeout: 5
  write_timeout: 10
  idle_timeout: 120

proxy:
  target_addr: "http://localhost:8000"
  max_idle_conns: 100
  dial_timeout: 10

load_balancing:
  enabled: false
  algorithm: "round_robin"
  backends: []

tls:
  enabled: false
  cert_file: ""
  key_file: ""

logging:
  level: "info"
  format: "json"

metrics:
  enabled: false
  address: ":9090"

rate_limiting:
  enabled: false
  requests_per_second: 100
  burst: 50

caching:
  enabled: false
  default_ttl: 300
  max_size_mb: 100
```

Remember to adjust these settings based on your specific requirements and the capabilities of your infrastructure.
