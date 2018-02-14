# Raspberry Pi Bartender

## Features
* Follows selected recipe to automatically pour drink
* Auto-generated [REST API](https://github.com/dhubler/bartend/blob/master/docs/api.md) and [model](https://github.com/dhubler/bartend/blob/master/docs/api.svg) documentation
* Access to pumps for priming, clearing pumps
* Uses IETF management standards for YANG and RESTCONF using [FreeCONF](https://github.com/freeconf) library
* Open Source (MIT)
* Single server for API, hardware access, business logic, and user interface
* Simple source code in Go.
* Mobile-first UI in Polymer 3.0
* Relatively inexpensive hardware
* `systemd` script to start when Pi starts

![Enclosure](https://github.com/dhubler/bartend/blob/master/docs/enclosure.jpg "Enclosure")


![Hardware](https://github.com/dhubler/bartend/blob/master/docs/hardware.jpg "Hardware")


![UI](https://github.com/dhubler/bartend/blob/master/docs/user-interface.gif "UI")

## Building
### Dependencies

These dependencies need to be in your PATH:

* Go 1.9 or newer [instructions](https://golang.org/dl)
* Go `dep` [instructions](https://golang.github.io/dep/docs/installation.html) 
* `yarn` [instructions](https://yarnpkg.com/lang/en/docs/install/)
* `make`

For graphical documentation generation:

* Graphviz [instructions](https://www.graphviz.org/download/)

### Building

This will build everything, including binary for Pi and your current workstation so you can develop.

```
make
```

## OS Support

This should work on all OSes.  If you port the `Makefile` to Windows, let me know and I'll include it with this project.

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
export YANGPATH=\
  ./etc/yang:\
  ./src/bartend/vendor/github.com/freeconf/gconf/yang
./bin/bartend -config ./etc/bartend.json 
```

## Running on Pi

After the build, there will be a `bartend.tgz`.  Copy that to Pi and untar in `/opt`.  You can install the `/opt/bartend/etc/service` into `systemd` for starting service with Pi.

## Limitations/Future

The pumps are very slow, a drink can take close to a minute to make.  I experimented with other pumps but they do not self-prime.  Ultimately I think a gravity-fed system might be a better design. Ideas welcomed.

You really shouldn't power-off Pi w/o shutting down gracefully or you risk corrupting SD card.  Switch to read-only filesystem or add a powerdown button.

Weighing scale. I started working on augmenting pumps with a scale that would allow you to pour liquids that were not attached to pumps. I never completed the integration of the hardware for the scale.

## Background

This is a functional project, a hobby and an example application of [FreeCONF](https://github.com/freeconf) libary.

## Support/Questions

You can contact if you have any comments for suggestions.

* email: douglas@hubler.us
* twitter: @dhubler


