package main

import (
	"bartend"
	"os"

	"github.com/c2stack/c2g/meta/yang"
	"github.com/c2stack/c2g/node"
	"github.com/c2stack/c2g/restconf"
)

var configFile = "bartend.cfg"

func main() {
	cfg, err := os.OpenFile(configFile, os.O_RDWR, os.ModeExclusive)
	defer cfg.Close()
	if err != nil {
		panic(err)
	}
	yangPath := yang.YangPath()
	m, err := yang.LoadModule(yangPath, "bartend")
	if err != nil {
		panic(err)
	}

	var app bartend.Bartend
	root := node.NewBrowser(m, bartend.Node(&app))
	if err := root.Root().InsertFrom(node.NewJsonReader(cfg).Node()).LastErr; err != nil {
		panic(err)
	}
	rc := restconf.NewService(yangPath, root)
	rc.Port = ":8080"
	rc.Listen()
}
