package main

import (
	"os"
	"github.com/c2g/node"
	"github.com/c2g/meta/yang"
	"github.com/c2g/restconf"
	"pi"
	"github.com/c2g/c2"
)

var configFile = "isaac.cfg"

func main() {
	cfg, err := os.OpenFile(configFile, os.O_RDWR, os.ModeExclusive)
	defer cfg.Close()
	if err != nil {
		panic(err)
	}
	m, err := yang.LoadModule(yang.InternalYang(), "isaac")
	if err != nil {
		panic(err)
	}

	c2.Debug.Printf("here")

	var root *node.Browser
	var app pi.Isaac
	var handler pi.ApiHandler
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

