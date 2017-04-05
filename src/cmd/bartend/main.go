package main

import (
	"bartend"
	"flag"
	"fmt"
	"os"

	"github.com/c2stack/c2g/conf"
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
	ypath := yang.YangPath()
	device := conf.NewLocalDeviceWithUi(ypath, web)
	device.Add("bartend", bartend.Node(app))

	rc := restconf.NewManagement(device)
	device.Add("restconf", restconf.Node(rc))

	if err := device.ApplyStartupConfig(*configFileName); err != nil {
		panic(err)
	}

	select {}
}
