# Diskmon

[![Build and Test](https://github.com/teleivo/diskmon/actions/workflows/build_test.yml/badge.svg)](https://github.com/teleivo/diskmon/actions/workflows/build_test.yml)
[![Release](https://img.shields.io/github/release/teleivo/diskmon.svg)](https://github.com/teleivo/diskmon/releases/latest)
[![Powered By: GoReleaser](https://img.shields.io/badge/powered%20by-goreleaser-green.svg)](https://github.com/goreleaser)

Diskmon will notify you via [Slack](https://slack.com) if a disk has reached a configurable size limit.

## Design Rationale

Diskmon was created to monitor volumes on [Digital Ocean](https://www.digitalocean.com/).
Digital ocean does not provide that feature in its droplet metrics (as of 09/2021).

We chose to implement this ourselves in the simplest way we could think of. We
did not want to setup [Prometheus](https://prometheus.io/) for this. If you
have prometheus already please use the [node exporter](https://github.com/prometheus/node_exporter).
It provides mount point monitoring and much more!

## Get started

Dowload a [pre-built binary](https://github.com/teleivo/diskmon/releases) or build the binary yourself

```sh
go build -o /usr/local/bin/diskmon
```

Change the destination to a path of your choice and make sure it can be found
by your shell via the $PATH variable.

Or run it directly

```sh
go run main.go -basedir <directory> -limit 65
```

which will write usage reports to stdout.

### Notifications

Notifications can be sent to
* stdout (by default)
* or Slack using a [Slack App](https://api.slack.com/start/building)

For Slack please follow the Slack documentation on how to create a Slack App Bot.
You can also follow this YouTube tutorial [Golang Tutorial: Build a Slack Bot](https://youtu.be/n-7l-N541u0).

You will then need to pass the Slack Bot User OAuth Token and the channel ID to
the binary via CLI flags.

Prefer passing credentials for example like so

```sh
diskmon -basedir <directory> -limit 65 \
  -slackToken $SLACK_TOKEN \
  -slackChannel $SLACK_CHANNEL
```

so that the credentials are not in your shell history.

## Limitations

The diskmon is not a general purpose disk monitor. It is specifically designed
for the use case we had ([see Design Rationale](#design-rationale)).

If you have prometheus already please use the [node exporter](https://github.com/prometheus/node_exporter).
It provides mount point monitoring and much more!

* Notifications can only be sent to Slack using a [Slack App](https://api.slack.com/start/building)
or to stdout
You can of course adapt the code to send notifications anywhere :smile:
* It will not discover all mount points for you like the [node exporter](https://github.com/prometheus/node_exporter).
You can only provide one directory in which your mount points should be.
For our use case the volumes we wanted to monitor were all in one directory.
