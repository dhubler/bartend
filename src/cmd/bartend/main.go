package main

import (
	"os"
	"github.com/c2g/node"
	"github.com/c2g/meta/yang"
	"github.com/c2g/restconf"
	"bartend"
)

var configFile = "bartend.cfg"

func main() {
	cfg, err := os.OpenFile(configFile, os.O_RDWR, os.ModeExclusive)
	defer cfg.Close()
	if err != nil {
		panic(err)
	}
	m, err := yang.LoadModule(yang.InternalYang(), "bartend")
	if err != nil {
		panic(err)
	}

	var root *node.Browser
	var app bartend.Bartend
	var handler bartend.ApiHandler
	root = node.NewBrowser(m, func() node.Node {
		return handler.Manage(root, &app)
	})
	if err := root.Root().Selector().InsertFrom(node.NewJsonReader(cfg).Node()).LastErr; err != nil {
		panic(err)
	}
	rc := restconf.NewService(root)

	rc.Port = ":8080"

	rc.Listen()

	//pi.LightLed()
}

