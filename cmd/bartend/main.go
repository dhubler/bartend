package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/dhubler/bartend"

	"github.com/freeconf/restconf/device"

	"github.com/freeconf/restconf"
	"github.com/freeconf/yang/source"
)

var configFileName = flag.String("config", "", "Configuration file")
var webDir = flag.String("web", "web/dist", "path to web directory")

// Start Bartend application.
//  YANGPATH=./etc/yang ./bin/bartend -config ./etc/bartend.json
func main() {
	flag.Parse()

	if len(*configFileName) == 0 {
		fmt.Fprint(os.Stderr, "Required 'config' parameter missing\n")
		os.Exit(-1)
	}

	app := bartend.NewBartend()
	ypath := source.Path(os.Getenv("YANGPATH"))
	dev := device.New(ypath)

	svr := restconf.NewServer(dev)

	svr.RegisterWebApp(*webDir, "index.html", "web")
	dev.Add("bartend", bartend.Node(app))

	if err := dev.ApplyStartupConfigFile(*configFileName); err != nil {
		panic(err)
	}

	select {}
}
