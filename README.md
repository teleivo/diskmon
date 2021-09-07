# Diskmon

[![build and test](https://github.com/teleivo/diskmon/actions/workflows/build_test.yml/badge.svg)](https://github.com/teleivo/diskmon/actions/workflows/build_test.yml)

Diskmon will notify you if a disk has reached a configurable size limit.

## Design Rationale

Diskmon was created to monitor volumes on [Digital Ocean](https://www.digitalocean.com/).
Digital ocean does not provide that feature in its droplet metrics (as of 09/2021).

We chose to implement this ourselves in the simplest way we could think of. We
did not want to setup [Prometheus](https://prometheus.io/) for this. If you
have prometheus already please use the [node exporter](https://github.com/prometheus/node_exporter).
It provides mount point monitoring and much more!

## Get started

### Using Binary

Build the binary or run directly using

```sh
go run main.go -basedir /home
```

### Using Docker

Build the image

```sh
docker build -t dockermon .
```

Run the image

```sh
docker run --volume /home:/home:ro diskmon -basedir /hom
```

## Limitations

The diskmon is not a general purpose disk monitor. It is specifically designed
for the use case we had ([see Design Rationale](#design-rationale))

If you have prometheus already please use the [node exporter](https://github.com/prometheus/node_exporter).
It provides mount point monitoring and much more!
