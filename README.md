# Diskmon

[![Build and Test](https://github.com/teleivo/diskmon/actions/workflows/build_test.yml/badge.svg)](https://github.com/teleivo/diskmon/actions/workflows/build_test.yml)
[![golangci-lint](https://github.com/teleivo/diskmon/actions/workflows/golangci-lint.yml/badge.svg)](https://github.com/teleivo/diskmon/actions/workflows/golangci-lint.yml)
[![codecov](https://codecov.io/gh/teleivo/diskmon/branch/main/graph/badge.svg?token=1VFP7UVS4Z)](https://codecov.io/gh/teleivo/diskmon)
[![Release](https://img.shields.io/github/release/teleivo/diskmon.svg)](https://github.com/teleivo/diskmon/releases/latest)
[![Powered By: GoReleaser](https://img.shields.io/badge/powered%20by-goreleaser-green.svg)](https://github.com/goreleaser)

Diskmon will notify you via [Slack](https://slack.com) if a disk has reached
(>=) a configurable size limit.

## Design Rationale

Diskmon was created to monitor volumes on [Digital Ocean](https://www.digitalocean.com/).
Digital ocean does not provide that feature in its droplet metrics (as of 09/2021).

We chose to implement this ourselves in the simplest way we could think of. We
did not want to setup [Prometheus](https://prometheus.io/) for this. If you
have prometheus already please use the [node exporter](https://github.com/prometheus/node_exporter).
It provides mount point monitoring and much more!

Since volumes are mounted under one directory on Digital Ocean (at least in our
setup) we chose to just get the disk usage for all its child directories
(non-recursive). In contrast, `node_exporter` will look at disk usage of your
mount points and let you for example ignore filesystems or mount points you are
not interested in.

We make a system call to [statfs](https://man.archlinux.org/man/statfs.2) to
get filesystem statistics. You can try `stat -f go.mod` locally to see what
usage information we get.

## Get started

Dowload a [pre-built binary](https://github.com/teleivo/diskmon/releases).

Copy the binary to a path that is dicoverable by your shell via the $PATH
environment variable.

Run it

```sh
diskmon -basedir <directory> -limit 65
```

Run `diskmon --help` to see all the available flags and their defaults.

### Notifications

Notifications can be sent to
* stdout (by default)
* or Slack using a [Slack App](https://api.slack.com/start/building)

For Slack please follow the Slack documentation on how to create a Slack App Bot.
You can also follow this YouTube tutorial [Golang Tutorial: Build a Slack Bot](https://youtu.be/n-7l-N541u0).

You will then need to pass the Slack **Bot User OAuth Token** and the channel
ID to the binary via CLI flags.

Prefer passing credentials for example like so

```sh
diskmon -basedir <directory> -limit 65 \
  -slackToken $SLACK_TOKEN \
  -slackChannel $SLACK_CHANNEL
```

so that the credentials are not in your shell history.

### Running as a service

Read [Running diskmon as a service](./examples/README.md).

## Build from source

You need to have [Go 1.16](https://golang.org) installed to build the binary yourself.

To build you can run

```sh
go build -o diskmon
```

Or build and run diskmon directly

```sh
go run main.go -basedir <directory> -limit 65
```

to familiarize yourself with the flags. The above will write usage reports to
stdout.

## Limitations

The diskmon is not a general purpose disk monitor. It is specifically designed
for the use case we had ([see Design Rationale](#design-rationale)).

If you have prometheus already please use the [node exporter](https://github.com/prometheus/node_exporter).
It provides mount point monitoring and much more!

* Notifications can only be sent to Slack using a [Slack App](https://api.slack.com/start/building)
or to stdout.
You can of course adapt the code to send notifications anywhere :smile:
* It will not discover all mount points for you like the [node exporter](https://github.com/prometheus/node_exporter).
You can only provide one directory in which your mount points should be.
For our use case the volumes we wanted to monitor were all in one directory.
* Does not work on Windows due to the syscall we are making
* Only checks disk usage for non-privileged users. A privileged user might have
  more disk space available. This might not work for you as it depends on the
  user permissions of the user writing to the disk you want to monitor.
