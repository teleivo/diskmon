# Running diskmon as a service

[systemd](https://systemd.io/)
> systemd is a suite of basic building blocks for a Linux system. It provides a
> system and service manager that runs as PID 1 and starts the rest of the
> system.

I provide a simple systemd [service file](./diskmon.service) that demonstrates
how you might be able to run diskmon as a service using `systemd`. You will
need to adjust paths to diskmon and mount points and flags to fit your
environment and needs.

## Installation

Install the diskmon service

```sh
sudo cp ./examples/diskmon.service /etc/systemd/system/diskmon.service
sudo chmod 664 /etc/systemd/system/diskmon.service
```

Adapt the service file to your needs!

Start diskmon to see if systemd can start diskmon successfully

```sh
sudo systemctl start diskmon
```

Check the service status

```sh
sudo systemctl status diskmon
```

or the diskmon logs

```sh
sudo journalctl --follow --unit diskmon --boot
```

If all went well, enable diskmon to start automatically at boot

```sh
sudo systemctl enable diskmon
```

Reboot and check the status, logs again to see if all is well.

## Limitations

* have not yet figured out how to securly pass the Slack API token. You could
  incorporate a subshell command that for example calls your password manager
  to pass the token in. This way you would not add the token in plain text into
  the service file which is readable by others on your system.
