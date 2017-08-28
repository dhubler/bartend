package main

import (
	"bartend"
	"flag"
	"fmt"
	"os"

	"github.com/c2stack/c2g/device"

	"github.com/c2stack/c2g/meta"
	"github.com/c2stack/c2g/meta/yang"
	"github.com/c2stack/c2g/restconf"
)

var configFileName = flag.String("config", "", "Configuration file")

// Start Bartend application.
//  YANGPATH=./etc/yang:./etc/c2-yang ./bin/bartend -config ./etc/bartend.json
func main() {
	flag.Parse()
	if len(*configFileName) == 0 {
		fmt.Fprint(os.Stderr, "Required 'config' parameter missing\n")
		os.Exit(-1)
	}

	app := bartend.NewBartend()

	web := &meta.FileStreamSource{Root: "web"}
	d := device.NewWithUi(yang.YangPath(), web)
	d.Add("bartend", bartend.Node(app))

	restconf.NewServer(d)

	if err := d.ApplyStartupConfigFile(*configFileName); err != nil {
		panic(err)
	}

	select {}
}
