# Raspberry Pi Bartender

## Features
* Follows selected recipe to automatically pour drink
* Auto-generated [REST API](https://github.com/dhubler/bartend/blob/master/docs/api.md) documentation
* Access to pumps for priming, clearing pumps
* Uses IETF management standards for YANG and RESTCONF using [FreeCONF](https://github.com/freeconf) library
* Open Source (MIT)
* Single server for API, hardware access, business logic, and user interface
* Approximately 800 lines of Go and 450 lines of html/css/JS
* Unit tests
* Mobile-first UI
* Relatively inexpensive hardware
* `systemd` script to start when Pi starts

![Enclosure](https://github.com/dhubler/bartend/blob/master/docs/enclosure.jpg "Enclosure")


![Hardware](https://github.com/dhubler/bartend/blob/master/docs/hardware.jpg "Hardware")


![UI](https://github.com/dhubler/bartend/blob/master/docs/user-interface.gif "UI")

## Building
### Dependencies

These dependencies need to be in your PATH:

* Go 1.18 or newer [instructions](https://golang.org/dl)
* `npm` [instructions](https://docs.npmjs.com/downloading-and-installing-node-js-and-npm)
* `make`

### Building

This will build everything, including binary for Pi and your current workstation so you can develop.

```
make pi
```

## OS Support

This should work on all OSes.  For Windows use WSL2 to build.

## Hardware

* Raspberry Pi Model 3/ 2GB or bigger SD card w/Raspian. Bartender application is around 8MB total size.
* 12v power supply, 1.5 Amp or better for driving Pi and pumps.
* Wiring
* 1-8 peristalic pumps. [Link](http://a.co/heFuT9v)
* Hosing
* Custom case
* 12v to 5v converter (old cigarette lighter USB charger will work)

## Running on local workstation

```
make run
```

## Running on Pi

After the build, there will be a `bartend.tgz`.  Copy that to Pi and untar in `/opt`.  You can install the `/opt/bartend/etc/service` into `systemd` for starting service with Pi.

## Limitations/Future

The pumps are very slow, a drink can take close to a minute to make.  I experimented with other pumps but they do not self-prime.  Ultimately I think a gravity-fed system might be a better design. Ideas welcomed.

You really shouldn't power-off Pi w/o shutting down gracefully or you risk corrupting SD card.  Switch to read-only filesystem or add a powerdown button.

## Background

This is a functional project, a hobby and an example application of [FreeCONF](https://github.com/freeconf) libary.

## Support/Questions

You can contact me if you have any comments for suggestions.

* email: douglas@hubler.us
