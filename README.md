# go-common-tools
Common tools for Bottlenose projects written in Go.

## Logger
`go-common-tools/logger` includes basic functionality to format messages into the bunyan format. It is pretty self explanatory, especially for those familiar with bunyan. It will write logs to stdout by default, unless a file path is provided when the logger is initialized. Loggers can be created using `NewLogger()` or `NewBufferedLogger()` if buffered output is desired.

## Metrics
`go-common-tools/metrics` provides wrapping functionality around the official golang prometheus client: `github.com/prometheus/client_golang/prometheus`. Currently supported [metrics](http://prometheus.io/docs/concepts/metric_types/) include counters, counterVecs, and gauges. We can add as many metrics types as we'd like as we find uses for them.

## Config
`go-common-tools/config` provides a config file/environment variable configuration helper.

## TestHTTP
`go-common-tools/testhttp` provides MockHTTP.
