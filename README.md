# Jellyfin exporter

[![Test & Build](https://github.com/rebelmediausa/jellyfin_exporter/actions/workflows/test.yml/badge.svg?branch=master)](https://github.com/rebelmediausa/jellyfin_exporter/actions/workflows/test.yml)
[![Current Release](https://img.shields.io/github/v/release/rebelmediausa/jellyfin_exporter)](https://github.com/rebelmediausa/jellyfin_exporter/releases/latest)
[![Docker Pulls](https://img.shields.io/docker/pulls/rebelmediausa/jellyfin-exporter)](https://hub.docker.com/r/rebelmediausa/jellyfin-exporter)
[![Go Report Card](https://goreportcard.com/badge/github.com/rebelmediausa/jellyfin_exporter)](https://goreportcard.com/report/github.com/rebelmediausa/jellyfin_exporter)

Prometheus exporter for Jellyfin Media System metrics exposed
in Go with pluggable metric collectors.


## Installation and Usage

If you are new to Prometheus and `jellyfin_exporter` there is
a [simple step-by-step guide](https://docs.rebelmedia.io/guides/jellyfin/exporter).

The `jellyfin_exporter` listens on HTTP port 9594 by default.
See the `--help` output for more options.

The flag `--jellyfin.token` is required. You can generate an API
Key in the Jellyfin admin dashboard. 

### Ansible

Coming Soon!


### Docker

The `jellyfin_exporter` is designed to monitor your Jellyfin Media System.

For situations where containerized deployment is needed, you will
need to set the Jellyfin URL flag to use the docker container hostname.

```bash
docker run -d \
  -p 9594:9594 \
  rebelmediausa/jellyfin-exporter:latest \
  --jellyfin.address=http://jellyfin:8096 \
  --jellyfin.token=TOKEN
```

For Docker compose, similar flag changes are needed.

```yaml
---
services:
  jellyfin_exporter:
    image: rebelmediausa/jellyfin-exporter:latest
    container_name: jellyfin_exporter
    command:
      - '--jellyfin.address=http://jellyfin:8096'
      - '--jellyfin.token=TOKEN'
    ports:
      - 9594:9594
    restart: unless-stopped
```


## Collectors

There is varying support for collectors.
The tables below list all existing collectors.

Collectors are enabled by providing a `--collector.<name>` flag.
Collectors that are enabled by default can be disabled
by providing a `--no-collector.<name>` flag.
To enable only some specific collector(s),
use `--collector.disable-defaults --collector.<name> ...`.


### Enabled by default

| Name    | Description                                        |
|---------|----------------------------------------------------|
| media   | Exposes media totals in the system by type.        |
| playing | Exposes media that users are now playing.          |
| system  | Exposes if the Jellyfin server is online or not.   |
| users   | Exposes users and if they are currently connected. |


### Filtering enabled collectors

The `jellyfin_exporter` will expose all metrics from enabled collectors
by default. This is the recommended way to collect metrics to avoid errors.

For advanced use the `jellyfin_exporter` can be passed an optional list
of collectors to filter metrics. The `collect[]` parameter may be used
multiple times. In Prometheus configuration you can use this syntax under
the [scrape config](https://prometheus.io/docs/prometheus/latest/configuration/configuration/#<scrape_config>).

```
  params:
    collect[]:
      - foo
      - bar
```

This can be useful for having different Prometheus servers collect
specific metrics from nodes.


## Development building and running

Prerequisites:

* [Go compiler](https://golang.org/dl/)
* RHEL/CentOS: `glibc-static` package.

Building:

    git clone https://github.com/rebelmediausa/jellyfin_exporter.git
    cd jellyfin_exporter
    make build
    ./jellyfin_exporter <flags>

To see all available configuration flags:

    ./jellyfin_exporter --help


## Running tests

    make test


## TLS endpoint

**EXPERIMENTAL**

The exporter supports TLS via a new web configuration file.

```console
./jellyfin_exporter --web.config.file=web-config.yml
```

See the [exporter-toolkit web-configuration](https://github.com/prometheus/exporter-toolkit/blob/master/docs/web-configuration.md) for more details.
