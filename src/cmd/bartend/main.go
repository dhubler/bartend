package main

import (
	"bartend"
	"flag"
	"fmt"
	"net/http"
	"os"

	"github.com/c2stack/c2g/meta/yang"
	"github.com/c2stack/c2g/node"
	"github.com/c2stack/c2g/restconf"
)

var configFileName = flag.String("config", "", "Configuration file")
var port = flag.String("port", ":80", "Port for http listener")

// Start Bartend application.
//  YANGPATH=./etc/yang ./bin/bartend -config ./etc/bartend.cfg -port :8080
func main() {
	flag.Parse()
	if len(*configFileName) == 0 {
		fmt.Fprint(os.Stderr, "Required 'config' parameter missing\n")
		os.Exit(-1)
	}

	cfg, err := os.OpenFile(*configFileName, os.O_RDWR, os.ModeExclusive)
	defer cfg.Close()
	if err != nil {
		panic(err)
	}
	yangPath := yang.YangPath()
	m, err := yang.LoadModule(yangPath, "bartend")
	if err != nil {
		panic(err)
	}

	app := bartend.NewBartend()
	root := node.NewBrowser(m, bartend.Node(app))
	if err := root.Root().InsertFrom(node.NewJsonReader(cfg).Node()).LastErr; err != nil {
		panic(err)
	}
	rc := restconf.NewManagement(yangPath, root)
	rc.Handle("/ui/", http.StripPrefix("/ui/", http.FileServer(http.Dir("web"))))
	rc.Handle("/", http.RedirectHandler("/ui/index.html", http.StatusFound))
	//rc.Handle("/ui/", stock.)
	options := rc.Options()
	options.Port = *port
	rc.ApplyOptions(options)
}
